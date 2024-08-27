/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package response

import "encoding/json"

type ResponseWrapper struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PayloadData struct {
	TxID       string `json:"tx_id"`
	OriginTxID string `json:"origin_tx_id"`
}

type APIResponse struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// ToJSON serializes the APIResponse struct to JSON
func (r *PayloadData) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the APIResponse struct
func (r *PayloadData) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

// ToJSON serializes the APIResponse struct to JSON
func (r *ResponseWrapper) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the APIResponse struct
func (r *ResponseWrapper) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
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
	return r.Payload == "" && r.Timestamp == 0 && r.Signature == ""
}
