// Package identifierallowlist maps external field names to trusted SQL identifiers.
package identifierallowlist

import (
	"errors"
	"strings"
)

var ErrIdentifierNotAllowed = errors.New("SQL identifier not allowed")

// Allowlist is immutable after construction.
type Allowlist struct{ values map[string]string }

// New builds a case-insensitive external-name mapping and rejects unsafe SQL identifiers.
func New(values map[string]string) (Allowlist, error) {
	copyValues := make(map[string]string, len(values))
	for external, identifier := range values {
		key := strings.ToLower(strings.TrimSpace(external))
		if key == "" || !validIdentifier(identifier) {
			return Allowlist{}, ErrIdentifierNotAllowed
		}
		copyValues[key] = identifier
	}
	return Allowlist{values: copyValues}, nil
}

func (a Allowlist) Resolve(external string) (string, error) {
	identifier, ok := a.values[strings.ToLower(strings.TrimSpace(external))]
	if !ok {
		return "", ErrIdentifierNotAllowed
	}
	return identifier, nil
}

func validIdentifier(value string) bool {
	parts := strings.Split(value, ".")
	for _, part := range parts {
		if part == "" {
			return false
		}
		for i, r := range part {
			if !(r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9') {
				return false
			}
		}
	}
	return true
}
