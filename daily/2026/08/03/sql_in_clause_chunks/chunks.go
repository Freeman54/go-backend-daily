// Package sql_in_clause_chunks plans batches that stay below a database parameter limit.
package sql_in_clause_chunks

import "fmt"

func Plan[T any](values []T, limit int) ([][]T, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("parameter limit must be positive: %d", limit)
	}
	chunks := make([][]T, 0, (len(values)+limit-1)/limit)
	for start := 0; start < len(values); start += limit {
		end := min(start+limit, len(values))
		chunks = append(chunks, values[start:end:end])
	}
	return chunks, nil
}
