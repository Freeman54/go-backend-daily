// Package error_disclosure_policy separates client-safe messages from internal error details.
package error_disclosure_policy

import (
	"errors"
	"fmt"
)

const fallback = "internal server error"

type safeError struct{ message string }

func (e safeError) Error() string { return e.message }

func Safe(message string) error { return safeError{message: message} }

func WrapForOperation(operation string, err error) error {
	return fmt.Errorf("%s: %w", operation, err)
}

func PublicMessage(err error) string {
	var safe safeError
	if errors.As(err, &safe) {
		return safe.message
	}
	return fallback
}
