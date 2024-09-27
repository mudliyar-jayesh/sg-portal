package util

import (
	"regexp"
)

// IsValidEmail checks if the given string is a valid email address
func IsValidEmail(email string) bool {
	// Basic regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// IsValidMobileNumber checks if the given string is a valid mobile number
func IsValidMobileNumber(mobile string) bool {
	// Assuming mobile number must be digits and between 10 and 15 characters long
	const mobileRegex = `^\d{10,15}$`
	re := regexp.MustCompile(mobileRegex)
	return re.MatchString(mobile)
}

