/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/response"
)

func (h *Handler) CompletionHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var resp response.APIResponse
	err = resp.FromJSON(body)
	if err != nil {
		http.Error(w, "Failed to marshal response, "+err.Error(), http.StatusInternalServerError)
		return
	}
	dataSigned := string(resp.Payload) + fmt.Sprintf("%d", resp.Timestamp)
	sigD, _ := base64.StdEncoding.DecodeString(resp.Signature)

	err = certs.VerifySignature(dataSigned, string(sigD), h.Config)
	if err != nil {
		http.Error(w, "Unable to verify signature", http.StatusInternalServerError)
		return
	}

	// TODO Deserialize Completion Payload if needed be
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body))
}
