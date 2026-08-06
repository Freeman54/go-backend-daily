// Package tracing_baggage_budget limits propagated tracing metadata.
package tracing_baggage_budget

// Apply keeps allowlisted entries in priority order within entry and byte budgets.
// Bytes are counted as the UTF-8 byte lengths of key plus value.
func Apply(input map[string]string, priority []string, maxEntries, maxBytes int) (map[string]string, int) {
	kept := make(map[string]string)
	seen := make(map[string]struct{})
	used := 0
	for _, key := range priority {
		if _, duplicate := seen[key]; duplicate {
			continue
		}
		seen[key] = struct{}{}
		value, ok := input[key]
		if !ok || len(kept) >= maxEntries || used+len(key)+len(value) > maxBytes {
			continue
		}
		kept[key] = value
		used += len(key) + len(value)
	}
	return kept, len(input) - len(kept)
}
