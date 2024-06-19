package request

import (
	"encoding/json"
	"fmt"

	client "github.com/caelisco/http-client"
	"github.com/intergreatme/remote-kyc-util/config"
	"github.com/intergreatme/remote-kyc-util/response"
)

type RequestPayload struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
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

func AllowlistAPI(payload RequestPayload, cnf config.Configuration) (response.APIResponse, response.ErrorResponse, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return response.APIResponse{}, response.ErrorResponse{}, err
	}

	uri := fmt.Sprintf("%sv2/allowlist/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "https://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Body.String())

	return response.APIResponse{}, response.ErrorResponse{}, nil
}
