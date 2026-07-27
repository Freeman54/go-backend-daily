package sql_filter_builder

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Dialect int

const (
	DialectPostgres Dialect = iota
	DialectMySQL
)

type Filter struct {
	Column string
	Value  any
}

var safeColumn = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Build 用参数占位符拼装等值过滤条件，值不会进入 SQL 文本。
func Build(base string, filters []Filter, dialect Dialect) (string, []any, error) {
	if len(filters) == 0 {
		return base, nil, nil
	}
	predicates := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))
	for index, filter := range filters {
		if !safeColumn.MatchString(filter.Column) {
			return "", nil, fmt.Errorf("unsafe column %q", filter.Column)
		}
		placeholder, err := bindVar(dialect, index+1)
		if err != nil {
			return "", nil, err
		}
		predicates = append(predicates, filter.Column+" = "+placeholder)
		args = append(args, filter.Value)
	}
	return base + " WHERE " + strings.Join(predicates, " AND "), args, nil
}

func bindVar(dialect Dialect, position int) (string, error) {
	switch dialect {
	case DialectPostgres:
		return fmt.Sprintf("$%d", position), nil
	case DialectMySQL:
		return "?", nil
	default:
		return "", errors.New("unsupported SQL dialect")
	}
}
