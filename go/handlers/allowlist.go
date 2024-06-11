package handlers

import (
	"encoding/json"
	"net/http"
)

type AllowlistRequest struct {
	OriginTxID      string  `json:"origin_tx_id"`               // The caller's transaction ID - required
	OrderNumber     string  `json:"order_number,omitempty"`     // The caller's order number if any (optional)
	FirstName       string  `json:"first_name"`                 // Required
	LastName        string  `json:"last_name"`                  // Required
	Mobile          string  `json:"mobile"`                     // Contactable mobile number
	Email           string  `json:"email"`                      // Contactable email address
	IDNumber        string  `json:"id_number,omitempty"`        // ID number (mutually exclusive with passport_number)
	PassportNumber  string  `json:"passport_number,omitempty"`  // Passport number (mutually exclusive with id_number)
	PassportCountry string  `json:"passport_country,omitempty"` // Mandatory if passport_number defined
	BuildingComplex string  `json:"building_complex,omitempty"` // Optional
	Line1           string  `json:"line1"`                      // Required
	Line2           string  `json:"line2,omitempty"`            // Optional
	Province        string  `json:"province"`                   // Required
	PostCode        string  `json:"post_code"`                  // Required
	Country         string  `json:"country"`                    // Required
	Latitude        float64 `json:"latitude,omitempty"`         // Optional
	Longitude       float64 `json:"longitude,omitempty"`        // Optional
	PlusCode        string  `json:"plus_code,omitempty"`        // Optional
}

type SignedPayload struct {
	Payload   AllowlistRequest `json:"payload"`
	Timestamp int64            `json:"timestamp"`
	Signature string           `json:"signature"`
}

type APIResponse struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type ErrorResponse struct {
	Valid    bool              `json:"valid"`
	UnderAge bool              `json:"under_age"`
	Errors   map[string]string `json:"errors"`
}

// The possible errors are:
// • TOO_SHORT: The provided identity number was too short
// • TOO_LONG: The provided identity number was too long
// • INVALID_CHARACTERS: The provided identity number contains non-numeric characters
// • REPEATED_CHARACTERS: The provided identity number contains only a single repeated digit
// • INVALID_BIRTH_DATE: The provided identity number encodes an invalid birth date
// • CHECK_DIGIT_MISMATCH: The check digit did not match
// • RSA_ID_ON_FOREIGN_TRACK: RSA identity number not allowed on foreign track

func AllowlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch and validate the request body
	var allowlistRequest AllowlistRequest
	if err := json.NewDecoder(r.Body).Decode(&allowlistRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate mandatory fields (add validation logic here)
	if err := validateAllowlistRequest(&allowlistRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sign the payload with the private key from the selfsign package
	signedPayload, err := signPayload(allowlistRequest)
	if err != nil {
		http.Error(w, "Failed to sign payload", http.StatusInternalServerError)
		return
	}

	// Make the POST request to the Allowlist API
	apiResponse, err := makeAllowlistAPICall(signedPayload)
	if err != nil {
		http.Error(w, "API call failed", http.StatusInternalServerError)
		return
	}

	// Handle and verify the API response
	if err := handleAPIResponse(apiResponse); err != nil {
		http.Error(w, "Failed to handle API response", http.StatusInternalServerError)
		return
	}

	// Store the response in the database
	if err := storeResponseInDB(apiResponse); err != nil {
		http.Error(w, "Failed to store response", http.StatusInternalServerError)
		return
	}

	// Respond with the API response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func validateAllowlistRequest(req *AllowlistRequest) error {
	// TODO
	// Implement validation logic
	// Check for mandatory fields and return an error if validation fails
	return nil
}

func signPayload(req AllowlistRequest) (SignedPayload, error) {
	// TODO
	// Implement the logic to sign the payload using the selfsign package
	// Return the signed payload or an error if signing fails
	return SignedPayload{}, nil
}

func makeAllowlistAPICall(payload SignedPayload) (APIResponse, error) {
	// TODO
	// Implement the logic to make the POST request to the Allowlist API
	// Return the API response or an error if the request fails
	return APIResponse{}, nil
}

func handleAPIResponse(response APIResponse) error {
	// TODO
	// Implement the logic to handle and verify the API response
	// Return an error if the response verification fails
	return nil
}

func storeResponseInDB(response APIResponse) error {
	// TODO
	// Implement the logic to store the API response in the database
	return nil
}

// Check if request is a POST else return method not allowed
// Fetch body
// Body should contain the allowlist JSON Schema
// With madatory fields required firstname, lastname, id number or passport number if foreign, then contact number and address (address is kinda optional)

// {
// 	origin_tx_id                         the callers transaction id - required
// 	order_number                         the callers order number if any (optional)

// 	first_name
// 	last_name
// 	mobile                               contactable mobile number
// 	email                                contactable email address

// 	id_number                            id_number & passport_number are mutually exclusive
// 	passport_number
// 	passport_country                     mandatory if passport_number defined

// 	building_complex
// 	line1
// 	line2
// 	province
// 	post_code
// 	country
// 	latitude
// 	longitude
// 	plus_code                            https://plus.codes/
// }

// after validating the necessary fields, start constructing the actual body to hit the Allowlist API for integration
// sign the payload with private key extracted from the certificate provided (use the selfsign package)
// It will be a POST, the content-encoding must be gzip
// Content-type is application/json
// Construct the body:
// Generic API JSON Schema with:
// 	payload: Allowlist JSON Schema
// 	timestamp: the UTC milliseconds since the epoch.
// 	signature: Base64 RSA Signed SHA512 of payload + timestamp

// Do the call. URI i'll provide myself.

// an example payload would be:
// 	{
// 	"payload":
// 		"{\"origin_tx_id\":\"YOUR_UUID\",\"first_name\":\"Joe\",\"last_name\":\"Soap\",
// 		\"id_number\":\"8301015468183\",\"contact_number\":\"0825555001\",\"mobile\":\"0825555001\",\"order_number\":\"00001\",
// 		\"building_complex\":\"Unit 3 Melrose Arch\",\"line1\":\"72 4th Avenue\",\"line2\":\"\",
// 		\"province\":\"Gauteng\",\"post_code\":\"2196\",\"country\":\"South Africa\",
// 		\"latitude\":\"-26.1024693\",\"longitude\":\"28.058389\",\"plus_code\":\"V3J3+9F Sandton\",
// 	"timestamp": 1551767491,
// 	"signature":
// 		"MmFiMGQzNWY0ZDI4YzdmZWE0NmU1M2FkNWJhZmFhYzAxOWY0ZmEwYTgzZDc3ZjFkMmFjNjYxYzM3Y
// 		mViOWJlNjUyZGI4YTFkY2IwM2ViM2MzMmJmMmQyNmIzMWNmYzQ4OTA0YzQ4Y2I0Y2I0MTA1NmIy
// 		NjdlYjlkZDY3ZWE3ZTM="
// }

// 	Now we need to handle the response.

// 	If Response is ok. Then handle response.
// 	Example response:
// 	{
// 		"payload":  "{\"tx_id\":\"KEY_UUID\",\"origin_tx_id\":\"YOUR_UUID\"}",
// 		 "timestamp": 1551767491,
// 		 "signature":
// 		   "MmFiMGQzNWY0ZDI4YzdmZWE0NmU1M2FkNWJhZmFhYzAxOWY0ZmEwYTgzZDc3ZjFkMmFjNjYxYzM3Y
// 		   mViOWJlNjUyZGI4YTFkY2IwM2ViM2MzMmJmMmQyNmIzMWNmYzQ4OTA0YzQ4Y2I0Y2I0MTA1NmIy
// 		   NjdlYjlkZDY3ZWE3ZTM="
//    }

//    Now verify the signature use selfsign package verify function.

// 	If error, throw error.
// Error JSON payload example
// {
// 	"valid": Boolean
// 	"under_age": Boolean
// 	"errors": {
// 		“ERROR_IDENTIFIER”: “Description”
// 	}
// }
// 	The possible errors are:
// • TOO_SHORT: The provided identity number was too short
// • TOO_LONG: The provided identity number was too long
// • INVALID_CHARACTERS: The provided identity number contains non-numeric characters
// • REPEATED_CHARACTERS: The provided identity number contains only a single repeated digit
// • INVALID_BIRTH_DATE: The provided identity number encodes an invalid birth date
// • CHECK_DIGIT_MISMATCH: The check digit did not match
// • RSA_ID_ON_FOREIGN_TRACK: RSA identity number not allowed on foreign track

// now store the response in the table
