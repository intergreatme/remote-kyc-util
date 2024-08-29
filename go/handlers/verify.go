/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */
package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/request"
)

// API end point to verify if the signature used is valid.
func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	timestampStr := r.FormValue("timestamp")

	// Convert string to int64
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}

	a := request.RequestPayload{
		Payload:   r.FormValue("payload"),
		Timestamp: timestamp,
		Signature: r.FormValue("signature"),
	}

	dataSigned := string(a.Payload) + fmt.Sprintf("%d", a.Timestamp)

	// Signature is base64 encoded, decode in order to verify.
	sigD, _ := base64.StdEncoding.DecodeString(a.Signature)
	err = certs.VerifySignature(dataSigned, string(sigD), h.Config)
	if err != nil {
		http.Error(w, "Unable to verify signature", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Verified"))
}
