/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package validator

import "regexp"

// Precompile the regex patterns
var (
	rgxEmail    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	rgxPassport = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	rgxMobile   = regexp.MustCompile(`^\+?[0-9\s\-()]+$`)
)

// IsValidEmail checks if the given value is a valid email address
func IsValidEmail(value string) bool {
	if len(value) > 254 {
		return false
	}
	return rgxEmail.MatchString(value)
}

// IsValidPassport checks if the given value is a valid passport number
func IsValidPassport(passport string) bool {
	return rgxPassport.MatchString(passport)
}

// IsValidMobile checks if the given value is a valid mobile number
func IsValidMobile(mobile string) bool {
	return rgxMobile.MatchString(mobile)
}
