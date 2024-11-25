package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

const (
	RULE_OLDER_THAN  = "older_than"
	RULE_REQUIRED_IF = "required_if"
)

type ValidationError struct {
	Key          string
	ErrorMessage string
}

// ParseValidationError parses an error string and returns a ValidationError struct.
func ParseValidationError(errorString string) (*ValidationError, error) {
	// Define regular expressions to extract the key and the error message.
	keyRegex := regexp.MustCompile(`Key:\s*'([^']+)'`)
	errorMessageRegex := regexp.MustCompile(`Error:(.+)`)

	// Find matches
	keyMatch := keyRegex.FindStringSubmatch(errorString)
	errorMessageMatch := errorMessageRegex.FindStringSubmatch(errorString)

	// Ensure matches were found
	if len(keyMatch) < 2 || len(errorMessageMatch) < 2 {
		return nil, fmt.Errorf("failed to parse validation error")
	}

	// Create a new ValidationError struct and return it
	validationError := &ValidationError{
		Key:          keyMatch[1],
		ErrorMessage: errorMessageMatch[1],
	}
	return validationError, nil
}

func OlderThan(fl validator.FieldLevel) bool {
	date, _ := time.Parse("2006-01-02", fl.Field().String())
	return time.Now().AddDate(-17, 0, 0).After(date)
}

func RequiredIf(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	otherFieldName := fl.Param()

	// Get the value of the other field
	otherFieldValue := fl.Parent().FieldByName(otherFieldName)

	// Check if the other field is empty
	if otherFieldValue.IsZero() {
		return true
	}

	// Check if the field is empty
	if fieldValue == "" {
		return false
	}

	return true
}

func IsNPWPValid(npwp string) bool {
	for _, c := range npwp {
		if c != '.' && c != '-' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func IsMimeTypeValid(mimeType string, allowedMimes []string) bool {
	for _, allowedMime := range allowedMimes {
		if mimeType == allowedMime {
			return true
		}
	}
	return false
}
