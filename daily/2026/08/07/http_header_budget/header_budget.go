// Package httpheaderbudget limits request header metadata before forwarding it.
package httpheaderbudget

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrFieldBudget = errors.New("header field budget exceeded")
	ErrByteBudget  = errors.New("header byte budget exceeded")
)

// Validate rejects headers exceeding maxFields or maxBytes. Bytes include every
// field name and value; separators are deliberately excluded for a stable policy.
func Validate(headers http.Header, maxFields, maxBytes int) error {
	if maxFields < 0 || maxBytes < 0 {
		return errors.New("header budgets must not be negative")
	}
	if len(headers) > maxFields {
		return fmt.Errorf("%w: got %d, limit %d", ErrFieldBudget, len(headers), maxFields)
	}
	used := 0
	for name, values := range headers {
		used += len(name)
		for _, value := range values {
			used += len(value)
		}
	}
	if used > maxBytes {
		return fmt.Errorf("%w: got %d, limit %d", ErrByteBudget, used, maxBytes)
	}
	return nil
}
