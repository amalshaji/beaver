package utils

import (
	"fmt"
	"net/mail"
	"strings"
)

// Remove any leading or trailing whitespace
func SanitizeString(value string) string {
	return strings.Trim(value, " ")
}

func ValidateEmail(input string) error {
	_, err := mail.ParseAddress(input)
	return err
}

func ValidatePassword(input string) error {
	if len(input) < 6 {
		return fmt.Errorf("password must be atleast 6 chars long")
	}
	return nil
}
