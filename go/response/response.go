package response

import "encoding/json"

type ErrorResponse struct {
	Valid    bool                 `json:"valid"`
	UnderAge bool                 `json:"under_age"`
	Errors   map[ErrorCode]string `json:"errors"`
}

// Define a custom type for the error codes
type ErrorCode string

type PayloadData struct {
	TxID       string `json:"tx_id"`
	OriginTxID string `json:"origin_tx_id"`
}

type APIResponse struct {
	Payload   PayloadField `json:"payload"`
	Timestamp int64        `json:"timestamp"`
	Signature string       `json:"signature"`
}

type PayloadField struct {
	Data PayloadData
}

func (pf *PayloadField) UnmarshalJSON(data []byte) error {
	var temp struct {
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	return json.Unmarshal([]byte(temp.Payload), &pf.Data)
}

func (pf PayloadField) MarshalJSON() ([]byte, error) {
	payloadData, err := json.Marshal(pf.Data)
	if err != nil {
		return nil, err
	}
	temp := struct {
		Payload string `json:"payload"`
	}{
		Payload: string(payloadData),
	}
	return json.Marshal(temp)
}

// ToJSONBytes converts the PayloadField to JSON bytes
func (pf *PayloadField) ToJSONBytes() ([]byte, error) {
	return json.Marshal(pf.Data)
}

// Define the possible error values
const (
	TOO_SHORT               ErrorCode = "TOO_SHORT"               // The provided identity number was too short
	TOO_LONG                ErrorCode = "TOO_LONG"                // The provided identity number was too long
	INVALID_CHARACTERS      ErrorCode = "INVALID_CHARACTERS"      // The provided identity number contains non-numeric characters
	REPEATED_CHARACTERS     ErrorCode = "REPEATED_CHARACTERS"     // The provided identity number contains only a single repeated digit
	INVALID_BIRTH_DATE      ErrorCode = "INVALID_BIRTH_DATE"      // The provided identity number encodes an invalid birth date
	CHECK_DIGIT_MISMATCH    ErrorCode = "CHECK_DIGIT_MISMATCH"    // The check digit did not match
	RSA_ID_ON_FOREIGN_TRACK ErrorCode = "RSA_ID_ON_FOREIGN_TRACK" // RSA identity number not allowed on foreign track
)

// IsEmpty checks if the APIResponse struct is empty
func (r APIResponse) IsEmpty() bool {
	return r.Payload == PayloadField{} && r.Timestamp == 0 && r.Signature == ""
}

// IsEmpty checks if the ErrorResponse struct is empty
func (e ErrorResponse) IsEmpty() bool {
	return !e.Valid && !e.UnderAge && len(e.Errors) == 0
}
