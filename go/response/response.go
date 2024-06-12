package response

type APIResponse struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type ErrorResponse struct {
	Valid    bool                 `json:"valid"`
	UnderAge bool                 `json:"under_age"`
	Errors   map[ErrorCode]string `json:"errors"`
}

// Define a custom type for the error codes
type ErrorCode string

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
	return r.Payload == "" && r.Timestamp == 0 && r.Signature == ""
}

// IsEmpty checks if the ErrorResponse struct is empty
func (e ErrorResponse) IsEmpty() bool {
	return !e.Valid && !e.UnderAge && len(e.Errors) == 0
}
