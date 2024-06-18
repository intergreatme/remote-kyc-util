package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/allowlist"
	"github.com/intergreatme/remote-kyc-util/response"
)

type RequestPayload struct {
	Payload   allowlist.Allowlist `json:"payload"`
	Timestamp int64               `json:"timestamp"`
	Signature []byte              `json:"signature"`
}

// ToSignableBytes converts the RequestPayload to JSON bytes excluding the Signature field
func (r *RequestPayload) ToSignableBytes() ([]byte, error) {
	type requestPayloadForSign struct {
		Payload   allowlist.Allowlist `json:"payload"`
		Timestamp int64               `json:"timestamp"`
	}

	request := requestPayloadForSign{
		Payload:   r.Payload,
		Timestamp: r.Timestamp,
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
	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(b)
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, err
	}
	gzipWriter.Close()

	// Convert the buffer's content to a byte slice
	compressedData := buf.Bytes()

	// Prepare the HTTP request
	url := "http://kycfe:8080/KycFrontEndServices/api/integration/v2/allowlist/d76f8a8c-1fe7-444e-9503-91a4f5d8451f"
	// Make the HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(compressedData))
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, fmt.Errorf("failed to make API call: %v", err)
	}
	defer resp.Body.Close()
	log.Println("req: ", resp)
	// Read the response body
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
