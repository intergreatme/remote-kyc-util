/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package allowlist

import (
	"encoding/json"
	"errors"

	"github.com/intergreatme/remote-kyc-util/validator"
)/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */



type Allowlist struct {
	OriginTxID      string `json:"origin_tx_id"`               // The caller's transaction ID - required
	OrderNumber     string `json:"order_number,omitempty"`     // The caller's order number if any (optional)
	FirstName       string `json:"first_name"`                 // Required
	LastName        string `json:"last_name"`                  // Required
	Mobile          string `json:"mobile"`                     // Contactable mobile number
	Email           string `json:"email"`                      // Contactable email address
	IDNumber        string `json:"id_number,omitempty"`        // ID number (mutually exclusive with passport_number)
	PassportNumber  string `json:"passport_number,omitempty"`  // Passport number (mutually exclusive with id_number)
	PassportCountry string `json:"passport_country,omitempty"` // Mandatory if passport_number defined
	BuildingComplex string `json:"building_complex,omitempty"` // Optional
	Line1           string `json:"line1"`                      // Required
	Line2           string `json:"line2,omitempty"`            // Optional
	Province        string `json:"province"`                   // Required
	PostCode        string `json:"post_code"`                  // Required
	Country         string `json:"country"`                    // Required
}

func (a Allowlist) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Allowlist) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a Allowlist) Validate() error {
	// Implement validation logic
	if a.OriginTxID == "" {
		return errors.New("origin transaction ID is required")
	}
	if a.FirstName == "" {
		return errors.New("first name is required")
	}
	if a.LastName == "" {
		return errors.New("last name is required")
	}
	if a.Line1 == "" {
		return errors.New("address line 1 is required")
	}
	if a.Province == "" {
		return errors.New("province is required")
	}
	if a.PostCode == "" {
		return errors.New("post code is required")
	}
	if a.Country == "" {
		return errors.New("country is required")
	}
	if a.PassportNumber == "" && a.IDNumber == "" {
		return errors.New("id number or passport number is required")
	}

	if a.Email != "" && !validator.IsValidEmail(a.Email) {
		return errors.New("invalid email address")
	}

	if a.Mobile != "" && !validator.IsValidMobile(a.Mobile) {
		return errors.New("invalid mobile number")
	}

	if a.PassportNumber != "" && !validator.IsValidPassport(a.PassportNumber) {
		return errors.New("invalid passport number")
	}

	if a.IDNumber != "" {
		resp := validator.IsValidID(a.IDNumber)
		if resp.HasError || !resp.ModulusCheck {
			return errors.New("invalid South African ID")
		}
	}

	return nil
}
