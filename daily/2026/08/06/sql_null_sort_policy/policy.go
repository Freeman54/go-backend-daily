// Package sql_null_sort_policy safely builds an ORDER BY clause.
package sql_null_sort_policy

import (
	"errors"
	"regexp"
)

type Direction string
type NullPlacement string

const (
	Ascending    Direction     = "ASC"
	Descending   Direction     = "DESC"
	NullsDefault NullPlacement = ""
	NullsFirst   NullPlacement = "NULLS FIRST"
	NullsLast    NullPlacement = "NULLS LAST"
)

var identifier = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.]*$`)

// BuildOrderBy maps an external sort key to a trusted database column.
func BuildOrderBy(key string, direction Direction, nulls NullPlacement, allowed map[string]string) (string, error) {
	column, ok := allowed[key]
	if !ok || !identifier.MatchString(column) {
		return "", errors.New("unknown or unsafe sort column")
	}
	if direction != Ascending && direction != Descending {
		return "", errors.New("invalid direction")
	}
	if nulls != NullsDefault && nulls != NullsFirst && nulls != NullsLast {
		return "", errors.New("invalid null placement")
	}
	if nulls == NullsDefault {
		return "ORDER BY " + column + " " + string(direction), nil
	}
	return "ORDER BY " + column + " " + string(direction) + " " + string(nulls), nil
}
