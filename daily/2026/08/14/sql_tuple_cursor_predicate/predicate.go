package tuplecursor

import (
	"fmt"
	"strings"
)

// After 构造同一排序方向下的复合 keyset 游标谓词。
func After(columns []string, values []any, descending bool) (string, []any, error) {
	if len(columns) == 0 || len(columns) != len(values) {
		return "", nil, fmt.Errorf("列和值必须非空且数量一致")
	}
	for _, column := range columns {
		if !validIdentifier(column) {
			return "", nil, fmt.Errorf("非法列名 %q", column)
		}
	}

	op := ">"
	if descending {
		op = "<"
	}
	parts := make([]string, len(columns))
	args := make([]any, 0, len(columns)*(len(columns)+1)/2)
	for i := range columns {
		terms := make([]string, 0, i+1)
		for j := 0; j < i; j++ {
			terms = append(terms, fmt.Sprintf("%s = ?", columns[j]))
			args = append(args, values[j])
		}
		terms = append(terms, fmt.Sprintf("%s %s ?", columns[i], op))
		args = append(args, values[i])
		parts[i] = "(" + strings.Join(terms, " AND ") + ")"
	}
	return "(" + strings.Join(parts, " OR ") + ")", args, nil
}

func validIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}
