// Package orderallowlist safely builds an ORDER BY clause from API fields.
package orderallowlist

import (
	"fmt"
	"strings"
)

type Direction string

const (
	Asc  Direction = "ASC"
	Desc Direction = "DESC"
)

type Field struct {
	Name      string
	Direction Direction
}

func Build(fields []Field, allowed map[string]string) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		column, ok := allowed[field.Name]
		if !ok {
			return "", fmt.Errorf("unsupported sort field %q", field.Name)
		}
		direction := field.Direction
		if direction == "" {
			direction = Asc
		}
		if direction != Asc && direction != Desc {
			return "", fmt.Errorf("unsupported direction %q", direction)
		}
		parts = append(parts, column+" "+string(direction))
	}
	return "ORDER BY " + strings.Join(parts, ", "), nil
}
