/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */
package update

import "encoding/json"

type Update struct {
	TxID           string `json:"tx_id"`
	IDNumber       string `json:"id_number,omitempty"`       // ID number (mutually exclusive with passport_number)
	PassportNumber string `json:"passport_number,omitempty"` // Passport number (mutually exclusive with id_number)
	FirstName      string `json:"first_name"`                // Required
	LastName       string `json:"last_name"`                 // Required
	ContactNumber  string `json:"contact_number"`
}

func (a Update) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Update) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}
