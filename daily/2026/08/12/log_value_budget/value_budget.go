// Package valuebudget bounds structured log field values while preserving valid UTF-8.
package valuebudget

import "unicode/utf8"

const marker = "…"

// Truncate returns value unchanged when it fits, otherwise a UTF-8-safe prefix and marker.
func Truncate(value string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(value) <= maxBytes {
		return value
	}
	if maxBytes < len(marker) {
		return ""
	}
	end := maxBytes - len(marker)
	for end > 0 && !utf8.ValidString(value[:end]) {
		end--
	}
	return value[:end] + marker
}

// Fields applies the same byte budget to every value and returns a detached map.
func Fields(fields map[string]string, maxBytes int) map[string]string {
	result := make(map[string]string, len(fields))
	for key, value := range fields {
		result[key] = Truncate(value, maxBytes)
	}
	return result
}
