package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/intergreatme/remote-kyc-util/allowlist"
	"github.com/intergreatme/remote-kyc-util/request"
	"github.com/intergreatme/selfsign"
)

func (h *Handler) AllowlistHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch and validate the request body
	var a allowlist.Allowlist
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a new origin_tx_id and order number
	a.OriginTxID = uuid.New().String()
	a.OrderNumber = fmt.Sprintf("%d%04d", time.Now().Unix(), rand.Intn(10000))

	// Validate mandatory fields (add validation logic here)
	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Serialize the payload to JSON
	payloadJSON, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	privateKey, err := selfsign.LoadPrivateKey("/keys/key.pem", h.Config.PvtKeyPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sign the serialized payload
	signature, err := selfsign.SignData(privateKey, payloadJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestBody := request.RequestPayload{
		Payload:   a,
		Timestamp: time.Now().Unix(),
		Signature: signature,
	}

	// Make the POST request to the Allowlist API
	respOK, respErr, err := request.AllowlistAPI(requestBody, h.Config.ID)
	if err != nil {
		http.Error(w, "API call failed", http.StatusInternalServerError)
		return
	}
	log.Println(respErr)

	if !respOK.IsEmpty() {
		cert, err := selfsign.LoadPEMCertificate("/keys/igm_certs.pem")
		if err != nil {
			http.Error(w, "Unable to load certificate", http.StatusInternalServerError)
			return
		}

		// Convert the payload to JSON bytes
		payloadBytes, err := respOK.Payload.ToJSONBytes()
		if err != nil {
			http.Error(w, "Unable to convert payload to JSON bytes", http.StatusInternalServerError)
			return
		}

		// Handle and verify the API response
		err = selfsign.VerifySignature(cert.PublicKey.(*rsa.PublicKey), payloadBytes, []byte(respOK.Signature))
		if err != nil {
			http.Error(w, "Unable to verify signature", http.StatusInternalServerError)
			return
		}

		// Store the OK response in the database
		sql := `INSERT INTO transactions 
		(origin_tx_id, tx_id, order_number, company, config_id, payload, response) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = h.DB.Exec(sql, respOK.Payload.Data.OriginTxID, respOK.Payload.Data.TxID, a.OrderNumber, "IGM-Test", h.Config.ID, payloadJSON, respOK)
		if err != nil {
			http.Error(w, "Unable to insert new record into transaction table", http.StatusInternalServerError)
			return
		}
		log.Println("Ok Response stored in db.")
	}
	if !respErr.IsEmpty() {

		// Store the ERROR response in the database
		sql := `INSERT INTO transactions 
	(origin_tx_id, tx_id, order_number, company, config_id, payload, response, errors) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = h.DB.Exec(sql, respOK.Payload.Data.OriginTxID, respOK.Payload.Data.TxID, a.OrderNumber, "IGM-Test", h.Config.ID, payloadJSON, respErr, respErr.Errors)
		if err != nil {
			http.Error(w, "Unable to insert new record into transaction table", http.StatusInternalServerError)
			return
		}
		log.Println("Error Response stored in db.")

	}
	// Respond with the API response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
}

// Check if request is a POST else return method not allowed
// Fetch body
// Body should contain the allowlist JSON Schema
// With madatory fields required firstname, lastname, id number or passport number if foreign, then contact number and address (address is kinda optional)

// {
// 	origin_tx_id                         the callers transaction id - required
// 	order_number                         the callers order number if any (optional)

// 	first_name
// 	last_name
// 	mobile                               contactable mobile number
// 	email                                contactable email address

// 	id_number                            id_number & passport_number are mutually exclusive
// 	passport_number
// 	passport_country                     mandatory if passport_number defined

// 	building_complex
// 	line1
// 	line2
// 	province
// 	post_code
// 	country
// 	latitude
// 	longitude
// 	plus_code                            https://plus.codes/
// }

// after validating the necessary fields, start constructing the actual body to hit the Allowlist API for integration
// sign the payload with private key extracted from the certificate provided (use the selfsign package)
// It will be a POST, the content-encoding must be gzip
// Content-type is application/json
// Construct the body:
// Generic API JSON Schema with:
// 	payload: Allowlist JSON Schema
// 	timestamp: the UTC milliseconds since the epoch.
// 	signature: Base64 RSA Signed SHA512 of payload + timestamp

// Do the call. URI i'll provide myself.

// an example payload would be:
// 	{
// 	"payload":
// 		"{\"origin_tx_id\":\"YOUR_UUID\",\"first_name\":\"Joe\",\"last_name\":\"Soap\",
// 		\"id_number\":\"8301015468183\",\"contact_number\":\"0825555001\",\"mobile\":\"0825555001\",\"order_number\":\"00001\",
// 		\"building_complex\":\"Unit 3 Melrose Arch\",\"line1\":\"72 4th Avenue\",\"line2\":\"\",
// 		\"province\":\"Gauteng\",\"post_code\":\"2196\",\"country\":\"South Africa\",
// 		\"latitude\":\"-26.1024693\",\"longitude\":\"28.058389\",\"plus_code\":\"V3J3+9F Sandton\",
// 	"timestamp": 1551767491,
// 	"signature":
// 		"MmFiMGQzNWY0ZDI4YzdmZWE0NmU1M2FkNWJhZmFhYzAxOWY0ZmEwYTgzZDc3ZjFkMmFjNjYxYzM3Y
// 		mViOWJlNjUyZGI4YTFkY2IwM2ViM2MzMmJmMmQyNmIzMWNmYzQ4OTA0YzQ4Y2I0Y2I0MTA1NmIy
// 		NjdlYjlkZDY3ZWE3ZTM="
// }

// 	Now we need to handle the response.

// 	If Response is ok. Then handle response.
// 	Example response:
// 	{
// 		"payload":  "{\"tx_id\":\"KEY_UUID\",\"origin_tx_id\":\"YOUR_UUID\"}",
// 		 "timestamp": 1551767491,
// 		 "signature":
// 		   "MmFiMGQzNWY0ZDI4YzdmZWE0NmU1M2FkNWJhZmFhYzAxOWY0ZmEwYTgzZDc3ZjFkMmFjNjYxYzM3Y
// 		   mViOWJlNjUyZGI4YTFkY2IwM2ViM2MzMmJmMmQyNmIzMWNmYzQ4OTA0YzQ4Y2I0Y2I0MTA1NmIy
// 		   NjdlYjlkZDY3ZWE3ZTM="
//    }

//    Now verify the signature use selfsign package verify function.

// 	If error, throw error.
// Error JSON payload example
// {
// 	"valid": Boolean
// 	"under_age": Boolean
// 	"errors": {
// 		“ERROR_IDENTIFIER”: “Description”
// 	}
// }
// 	The possible errors are:
// • TOO_SHORT: The provided identity number was too short
// • TOO_LONG: The provided identity number was too long
// • INVALID_CHARACTERS: The provided identity number contains non-numeric characters
// • REPEATED_CHARACTERS: The provided identity number contains only a single repeated digit
// • INVALID_BIRTH_DATE: The provided identity number encodes an invalid birth date
// • CHECK_DIGIT_MISMATCH: The check digit did not match
// • RSA_ID_ON_FOREIGN_TRACK: RSA identity number not allowed on foreign track

// now store the response in the table
