package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: " test ", want: "test"},
		{input: "test ", want: "test"},
		{input: " test", want: "test"},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.want, SanitizeString(tc.input))
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{input: " test ", hasError: true},
		{input: "test@beaver.com", hasError: false},
	}

	for _, tc := range tests {
		if tc.hasError {
			assert.Error(t, ValidateEmail(tc.input))
		} else {
			assert.NoError(t, ValidateEmail(tc.input))
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{input: "test", want: ErrPasswordNotLongEnough},
		{input: "test@beaver.com", want: nil},
	}

	for _, tc := range tests {
		assert.ErrorIs(t, ValidatePassword(tc.input), tc.want)
	}
}

func TestValidateSubdomain(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{input: "beaver-test-dev", want: nil},
		{input: "beaver.test", want: ErrInvalidSubdomain},
		{input: "-beaver", want: ErrInvalidSubdomain},
		{input: "beaver-", want: ErrInvalidSubdomain},
		{input: "beaver_test", want: ErrInvalidSubdomain},
	}

	for _, tc := range tests {
		assert.ErrorIs(t, ValidateSubdomain(tc.input), tc.want)
	}
}
