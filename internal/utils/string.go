package utils

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrPasswordNotLongEnough = errors.New("password must be atleast 6 chars long")

// Remove any leading or trailing whitespace
func SanitizeString(value string) string {
	return strings.TrimSpace(value)
}

func ValidateEmail(input string) error {
	_, err := mail.ParseAddress(input)
	return err
}

func ValidatePassword(input string) error {
	if len(input) < 6 {
		return ErrPasswordNotLongEnough
	}
	return nil
}
