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
	"github.com/intergreatme/remote-kyc-util/cancel"
	"github.com/intergreatme/remote-kyc-util/request"
)

func (h *Handler) CancelHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	a := cancel.Cancel{
		TxID:    r.FormValue("tx_id"),
		Comment: r.FormValue("comment"),
	}

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

	datasign := string(b) + fmt.Sprintf("%d", timestamp)

	signature, err := certcrypto.SignData(privateKey, []byte(datasign))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(signature)

	sig := base64.StdEncoding.EncodeToString(signature)
	reqBody := request.RequestPayload{
		Payload:   string(b),
		Timestamp: timestamp,
		Signature: sig,
	}

	// Make the POST request to the Allowlist API
	resp, err := request.CancelAPI(reqBody, h.Config)
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
	log.Printf("\n\nCancel API response: %s", payloadBody)

	// Add any extra handling here
	w.Write([]byte(payloadBody))

}
