package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/response"
)

type RequestPayload struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature []byte `json:"signature"`
}

// ToSignableBytes converts the RequestPayload to JSON bytes excluding the Signature field
func (r *RequestPayload) ToSignableBytes() ([]byte, error) {
	type requestPayloadForSign struct {
		Payload   string `json:"payload"`
		Timestamp string `json:"timestamp"`
	}

	request := requestPayloadForSign{
		Payload:   r.Payload,
		Timestamp: fmt.Sprintf("%d", r.Timestamp),
	}

	return json.Marshal(request)
}

// ToJSON serializes the RequestPayload struct to JSON
func (r RequestPayload) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the RequestPayload struct
func (r *RequestPayload) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

func AllowlistAPI(rp RequestPayload, configID string) (response.APIResponse, response.ErrorResponse, error) {
	b, err := rp.ToJSON()
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, err
	}

	// Compress the JSON payload with gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(b)); err != nil {
		log.Fatalf("failed to compress data: %v", err)
	}
	if err := gz.Close(); err != nil {
		log.Fatalf("failed to close gzip writer: %v", err)
	}

	url := "http://kycfe:8080/KycFrontEndServices/api/integration/v2/allowlist/" + configID

	// Print JSON data for debugging
	fmt.Println("JSON Data:", string(b))

	// Create the HTTP POST request using http.Client
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	// Make the HTTP POST request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to make API call: %v", err)
	}
	defer resp.Body.Close()
	log.Println(resp)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, err
	}

	// Handle different response types based on status code
	if resp.StatusCode == http.StatusOK {
		// Successful response
		var apiResp response.APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return response.APIResponse{}, response.ErrorResponse{}, fmt.Errorf("failed to decode API response: %v", err)
		}
		return apiResp, response.ErrorResponse{}, nil
	} else {
		// Error response
		var errResp response.ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return response.APIResponse{}, response.ErrorResponse{}, fmt.Errorf("failed to decode error response: %v", err)
		}
		return response.APIResponse{}, errResp, nil
	}
}
