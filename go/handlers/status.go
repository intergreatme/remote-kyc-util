package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/response"
)

// StatusUpdate represents the structure of the status update JSON payload
type StatusUpdate struct {
	TxID          string        `json:"tx_id"`
	OriginTxID    string        `json:"origin_tx_id"`
	Status        string        `json:"status"`
	DocumentTypes []string      `json:"document_types,omitempty"`
	ProfileState  *ProfileState `json:"profile_state,omitempty"`
	Comment       string        `json:"comment,omitempty"`
}

// ProfileState represents the profile state in the status update
type ProfileState struct {
	Identity         string `json:"identity"`
	Address          string `json:"address"`
	ProofOfResidence string `json:"proof_of_residence"`
	Consent          string `json:"consent"`
	Liveliness       string `json:"liveliness"`
}

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var respOK response.APIResponse
	err = respOK.FromJSON(body)
	if err != nil {
		http.Error(w, "Failed to marshal response, "+err.Error(), http.StatusInternalServerError)
		return
	}
	dataSigned := string(respOK.Payload) + fmt.Sprintf("%d", respOK.Timestamp)
	sigD, _ := base64.StdEncoding.DecodeString(respOK.Signature)

	err = certs.VerifySignature(dataSigned, string(sigD), h.Config)
	if err != nil {
		http.Error(w, "Unable to verify signature", http.StatusInternalServerError)
		return
	}

	// TODO Deserialize Status Payload if needed be
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(body))
}
