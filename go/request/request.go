package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/allowlist"
	"github.com/intergreatme/remote-kyc-util/response"

	client "github.com/caelisco/http-client"
)

type RequestPayload struct {
	Payload   allowlist.Allowlist `json:"payload"`
	Timestamp int64               `json:"timestamp"`
	Signature []byte              `json:"signature"`
}

// ToJSON serializes the RequestPayload struct to JSON
func (r RequestPayload) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the RequestPayload struct
func (r *RequestPayload) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

func CallAPI(rp RequestPayload) (response.APIResponse, response.ErrorResponse, error) {

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

	var opt client.RequestOptions
	opt.AddHeader("Content-Type", "application/json")
	opt.AddHeader("Content-Encoding", "gzip")

	config := "" // TODO
	resp, err := client.Post("url"+config, compressedData, opt)
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, fmt.Errorf("failed to make API call: %v", err)
	}

	// Read the response body
	respBody, err := io.ReadAll(&resp.Body)
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
