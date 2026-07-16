package baggageallowlist

import "strings"

func Filter(header string, allowed map[string]struct{}, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}

	parts := strings.Split(header, ",")
	result := make([]string, 0, len(parts))
	size := 0
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if _, ok := allowed[key]; !ok {
			continue
		}

		entry := key + "=" + value
		nextSize := size + len(entry)
		if len(result) > 0 {
			nextSize++
		}
		if nextSize > maxBytes {
			break
		}
		result = append(result, entry)
		size = nextSize
	}
	return strings.Join(result, ",")
}
