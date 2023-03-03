package utils

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

var (
	ErrPasswordNotLongEnough = errors.New("password must be atleast 6 chars long")
	ErrInvalidSubdomain      = errors.New("subdomain must contain only a-z, 0-9 and `-`(not leading or trailing)")
)

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

func ValidateSubdomain(subdomain string) error {
	matched, err := regexp.Match(`^([a-z0-9]+([-][a-z0-9]+)*)$`, []byte(subdomain))
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidSubdomain
	}
	return nil
}
