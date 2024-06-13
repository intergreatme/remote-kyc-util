package response

import "encoding/json"

type PayloadData struct {
	TxID       string `json:"tx_id"`
	OriginTxID string `json:"origin_tx_id"`
}

type PayloadField struct {
	Data PayloadData `json:"data"`
}

type APIResponse struct {
	Payload   PayloadField `json:"payload"`
	Timestamp int64        `json:"timestamp"`
	Signature string       `json:"signature"`
}

// ToSignableBytes converts the APIResponse to JSON bytes excluding the Signature field
func (r *APIResponse) ToSignableBytes() ([]byte, error) {
	type apiResponseForSign struct {
		Payload   PayloadField `json:"payload"`
		Timestamp int64        `json:"timestamp"`
	}

	response := apiResponseForSign{
		Payload:   r.Payload,
		Timestamp: r.Timestamp,
	}

	return json.Marshal(response)
}

// ToJSON serializes the APIResponse struct to JSON
func (r *APIResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the APIResponse struct
func (r *APIResponse) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// IsEmpty checks if the APIResponse struct is empty
func (r APIResponse) IsEmpty() bool {
	return r.Payload == PayloadField{} && r.Timestamp == 0 && r.Signature == ""
}
