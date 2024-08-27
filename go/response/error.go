/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package response

import (
	"encoding/json"
)

// Define a custom type for the error codes
type ErrorCode string

type ErrorResponse struct {
	Valid    bool                 `json:"valid"`
	UnderAge bool                 `json:"under_age"`
	Errors   map[ErrorCode]string `json:"errors"`
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

// IsEmpty checks if the ErrorResponse struct is empty
func (e ErrorResponse) IsEmpty() bool {
	return !e.Valid && !e.UnderAge && len(e.Errors) == 0
}

// MarshalJSON customizes the JSON encoding for ErrorResponse
func (e ErrorResponse) toJSON() ([]byte, error) {
	type Alias ErrorResponse // Create an alias to avoid recursion
	return json.Marshal(&struct {
		Errors map[string]string `json:"errors"`
		Alias
	}{
		Errors: convertMapErrorCodeToString(e.Errors),
		Alias:  (Alias)(e),
	})
}

// UnmarshalJSON customizes the JSON decoding for ErrorResponse
func (e *ErrorResponse) FromJSON(data []byte) error {
	type Alias ErrorResponse // Create an alias to avoid recursion
	aux := &struct {
		Errors map[string]string `json:"errors"`
		Alias
	}{
		Alias: (Alias)(*e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Valid = aux.Valid
	e.UnderAge = aux.UnderAge
	e.Errors = convertMapStringToErrorCode(aux.Errors)
	return nil
}

// Helper function to convert map[ErrorCode]string to map[string]string
func convertMapErrorCodeToString(m map[ErrorCode]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[string(k)] = v
	}
	return result
}

// Helper function to convert map[string]string to map[ErrorCode]string
func convertMapStringToErrorCode(m map[string]string) map[ErrorCode]string {
	result := make(map[ErrorCode]string)
	for k, v := range m {
		result[ErrorCode(k)] = v
	}
	return result
}
