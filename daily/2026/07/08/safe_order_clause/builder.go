package safeorderclause

import (
	"fmt"
	"strings"
)

type Builder struct {
	columns map[string]string
}

func New(columns map[string]string) *Builder {
	copied := make(map[string]string, len(columns))
	for key, value := range columns {
		copied[key] = value
	}
	return &Builder{columns: copied}
}

func (b *Builder) Build(field string, descending bool) (string, error) {
	column, ok := b.columns[field]
	if !ok {
		return "", fmt.Errorf("unsupported sort field: %s", field)
	}

	direction := "ASC"
	if descending {
		direction = "DESC"
	}

	return strings.Join([]string{
		column + " " + direction,
		"id ASC",
	}, ", "), nil
}
