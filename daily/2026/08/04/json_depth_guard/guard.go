package json_depth_guard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	ErrTooDeep         = errors.New("json nesting exceeds maximum depth")
	ErrInvalidMaxDepth = errors.New("maximum depth must be positive")
)

// Validate verifies JSON syntax and rejects documents nested deeper than maxDepth.
func Validate(data []byte, maxDepth int) error {
	if maxDepth <= 0 {
		return ErrInvalidMaxDepth
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	depth := 0
	seen := false
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			if depth != 0 {
				return io.ErrUnexpectedEOF
			}
			break
		}
		if err != nil {
			return fmt.Errorf("decode JSON: %w", err)
		}
		seen = true
		if delimiter, ok := token.(json.Delim); ok {
			switch delimiter {
			case '{', '[':
				depth++
				if depth > maxDepth {
					return ErrTooDeep
				}
			case '}', ']':
				depth--
			}
		}
	}
	if !seen {
		return io.ErrUnexpectedEOF
	}
	return nil
}
