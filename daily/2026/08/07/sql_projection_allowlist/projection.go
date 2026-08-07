// Package sqlprojectionallowlist safely turns approved columns into a SELECT list.
package sqlprojectionallowlist

import (
	"errors"
	"fmt"
	"strings"
)

// Select quotes requested columns after checking that each is allowed and unique.
func Select(columns []string, allowed map[string]struct{}) (string, error) {
	if len(columns) == 0 {
		return "", errors.New("at least one projected column is required")
	}
	seen := make(map[string]struct{}, len(columns))
	quoted := make([]string, 0, len(columns))
	for _, column := range columns {
		if _, ok := allowed[column]; !ok {
			return "", fmt.Errorf("column %q is not allowed", column)
		}
		if _, duplicate := seen[column]; duplicate {
			return "", fmt.Errorf("column %q is repeated", column)
		}
		seen[column] = struct{}{}
		quoted = append(quoted, `"`+strings.ReplaceAll(column, `"`, `""`)+`"`)
	}
	return strings.Join(quoted, ", "), nil
}
