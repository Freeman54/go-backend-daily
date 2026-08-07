// Package messagepartitionkey validates keys used to preserve message order.
package messagepartitionkey

import (
	"errors"
	"fmt"
)

// Validate accepts a non-empty, bounded ASCII partition key. The allowed
// punctuation is intentionally small so producers in different languages agree.
func Validate(key string, maxLength int) error {
	if maxLength < 1 {
		return errors.New("partition key maximum length must be positive")
	}
	if key == "" {
		return errors.New("partition key is required")
	}
	if len(key) > maxLength {
		return fmt.Errorf("partition key length %d exceeds limit %d", len(key), maxLength)
	}
	for _, char := range key {
		if !allowed(char) {
			return fmt.Errorf("partition key contains invalid character %q", char)
		}
	}
	return nil
}

func allowed(char rune) bool {
	return char >= 'a' && char <= 'z' ||
		char >= 'A' && char <= 'Z' ||
		char >= '0' && char <= '9' ||
		char == '-' || char == '_' || char == '.' || char == ':'
}
