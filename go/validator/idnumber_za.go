/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package validator

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type IDResult struct {
	DateOfBirth  time.Time
	Age          int
	Over18       bool
	Gender       string
	ModulusCheck bool
	HasError     bool
}

func IsValidID(idNumber string) IDResult {
	idNumber = strings.TrimSpace(idNumber)

	if len(idNumber) != 13 {
		log.Print("Invalid ID number: length must be 13 digits.")
		return IDResult{HasError: true}
	}

	if !isNumeric(idNumber) {
		log.Print("Invalid ID number: must contain only digits.")
		return IDResult{HasError: true}
	}

	dateOfBirth, err := getDateOfBirth(idNumber)
	if err != nil {
		log.Print("Invalid ID number:", err)
		return IDResult{HasError: true}
	}

	age := getAge(dateOfBirth)
	isOver18 := age >= 18
	gender := getGender(idNumber)
	validModulus := validateModulus(idNumber)

	r := IDResult{
		DateOfBirth:  dateOfBirth,
		Age:          age,
		Over18:       isOver18,
		Gender:       gender,
		ModulusCheck: validModulus,
	}

	return r
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func getDateOfBirth(idNumber string) (time.Time, error) {
	year, _ := strconv.Atoi("19" + idNumber[0:2])
	month, _ := strconv.Atoi(idNumber[2:4])
	day, _ := strconv.Atoi(idNumber[4:6])

	dateOfBirth := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if dateOfBirth.After(time.Now()) {
		return time.Time{}, fmt.Errorf("date of birth is in the future")
	}

	return dateOfBirth, nil
}

func getAge(dateOfBirth time.Time) int {
	now := time.Now()
	years := now.Year() - dateOfBirth.Year()

	if now.Month() < dateOfBirth.Month() || (now.Month() == dateOfBirth.Month() && now.Day() < dateOfBirth.Day()) {
		years--
	}

	return years
}

func getGender(idNumber string) string {
	genderDigit, _ := strconv.Atoi(idNumber[6:7])
	if genderDigit >= 5 {
		return "Male"
	}
	return "Female"
}

func validateModulus(idNumber string) bool {
	oddSum := 0
	evenSum := 0

	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(idNumber[i]))
		if i%2 == 0 {
			oddSum += digit
		} else {
			evenSum += digit * 2
			if digit >= 5 {
				evenSum -= 9
			}
		}
	}

	totalSum := oddSum + evenSum
	checkDigit := 10 - (totalSum % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}

	return checkDigit == int(idNumber[12]-'0')
}
