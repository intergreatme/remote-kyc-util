/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/intergreatme/certcrypto"
	"github.com/intergreatme/remote-kyc-util/allowlist"
	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/request"
	"github.com/intergreatme/remote-kyc-util/response"
)

func (h *Handler) AllowlistHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	a := allowlist.Allowlist{
		OriginTxID:      r.FormValue("origin_tx_id"),
		OrderNumber:     r.FormValue("order_number"),
		FirstName:       r.FormValue("first_name"),
		LastName:        r.FormValue("last_name"),
		Mobile:          r.FormValue("mobile"),
		Email:           r.FormValue("email"),
		IDNumber:        r.FormValue("id_number"),
		PassportNumber:  r.FormValue("passport_number"),
		PassportCountry: r.FormValue("passport_country"),
		BuildingComplex: r.FormValue("building_complex"),
		Line1:           r.FormValue("line1"),
		Line2:           r.FormValue("line2"),
		Province:        r.FormValue("province"),
		PostCode:        r.FormValue("post_code"),
		Country:         r.FormValue("country"),
	}

	// Validate mandatory fields (add validation logic here)
	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// NB: Timestamp modification required
	timestamp := time.Now().Unix()

	b, err := a.ToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pfxFile := filepath.Join(h.Config.CertDir, h.Config.PFXFilename)
	privateKey, _, err := certcrypto.ReadPKCS12(pfxFile, h.Config.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sign the serialized payload
	datasign := string(b) + fmt.Sprintf("%d", timestamp)

	signature, err := certcrypto.SignData(privateKey, []byte(datasign))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sig := base64.StdEncoding.EncodeToString(signature)
	reqBody := request.RequestPayload{
		Payload:   string(b),
		Timestamp: timestamp,
		Signature: sig,
	}

	// Make the POST request to the Allowlist API
	resp, err := request.AllowlistAPI(reqBody, h.Config)
	if err != nil {
		http.Error(w, "API call failed, "+err.Error(), http.StatusInternalServerError)
		return
	}

	payloadBody, err := io.ReadAll(&resp.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("\n\nAPI response: %s", payloadBody)

	// If error, fail fast.
	if resp.StatusCode != 200 {
		var rw response.ResponseWrapper
		// Errors are sent with the response wrapper, which required further deserializing
		err = rw.FromJSON(payloadBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("\n\nPayload response: %v", rw)

		// Uncomment this if you want to handle specific custom error responses like status code 406.
		// {
		//     "valid": false,
		//     "under_age": false,
		//     "errors": {
		//         “TOO_SHORT”: “The provided identity number was too short”,
		//         “INVALID_CHARACTERS”: “The provided identity number contains non-numeric characters”
		//     }
		// }
		// var respM response.ErrorResponse
		// err = respM.FromJSON([]byte(rw.Message))
		// if err != nil {
		// 	http.Error(w, "Unable to unmarshal response message", http.StatusInternalServerError)
		// 	return
		// }

		sql := `INSERT INTO transactions
		(origin_tx_id, tx_id, order_number, company, config_id, payload, response, errors)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = h.DB.Exec(sql, a.OriginTxID, nil, a.OrderNumber, "IGM-Test", h.Config.CompanyID, string(b), payloadBody, rw.Message)
		if err != nil {
			http.Error(w, "Unable to insert new record into transaction table", http.StatusInternalServerError)
			return
		}
		log.Printf("Error response stored in DB, error received: %d - %s", rw.Code, rw.Message)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payloadBody))
		return
	}

	// Happy route, new transaction created.
	var respOK response.APIResponse
	err = respOK.FromJSON(payloadBody)
	if err != nil {
		http.Error(w, "Failed to marshal response, "+err.Error(), http.StatusInternalServerError)
		return
	}

	dataSigned := string(respOK.Payload) + fmt.Sprintf("%d", respOK.Timestamp)

	// Signature is base64 encoded, decode in order to verify.
	sigD, _ := base64.StdEncoding.DecodeString(respOK.Signature)
	err = certs.VerifySignature(dataSigned, string(sigD), h.Config)
	if err != nil {
		http.Error(w, "Unable to verify signature", http.StatusInternalServerError)
		return
	}

	// Fetch TX ID from body
	var pd response.PayloadData
	err = pd.FromJSON([]byte(respOK.Payload))
	if err != nil {
		http.Error(w, "Unable to unmarshal response payload", http.StatusInternalServerError)
		return
	}

	log.Printf("OriginTxID: %s, TxID: %s", pd.OriginTxID, pd.TxID)
	sql := `INSERT INTO transactions
			(origin_tx_id, tx_id, order_number, company, config_id, payload, response)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = h.DB.Exec(sql, pd.OriginTxID, pd.TxID, a.OrderNumber, "IGM-Test", h.Config.CompanyID, string(b), string(payloadBody))
	if err != nil {
		http.Error(w, "Unable to insert new record into transaction table", http.StatusInternalServerError)
		return
	}
	log.Printf("Ok Response stored in db, for Origin TX ID: %s and TxID: %s", pd.OriginTxID, pd.TxID)

	// Respond with the API response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payloadBody))
}
