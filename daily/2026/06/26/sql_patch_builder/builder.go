package sqlpatchbuilder

import (
	"fmt"
	"sort"
	"strings"
)

type Patch map[string]any

type Builder struct {
	table   string
	allowed map[string]string
}

func NewBuilder(table string) Builder {
	return Builder{
		table:   table,
		allowed: make(map[string]string),
	}
}

func (b Builder) Allow(column, field string) Builder {
	b.allowed[field] = column
	return b
}

func (b Builder) BuildUpdate(patch Patch, where string, whereArgs ...any) (string, []any, error) {
	if len(patch) == 0 {
		return "", nil, fmt.Errorf("empty patch")
	}

	type pair struct {
		column string
		value  any
	}

	pairs := make([]pair, 0, len(patch))
	for field, value := range patch {
		column, ok := b.allowed[field]
		if !ok {
			return "", nil, fmt.Errorf("field %q is not allowed", field)
		}
		pairs = append(pairs, pair{column: column, value: value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].column < pairs[j].column
	})

	setClauses := make([]string, 0, len(pairs))
	args := make([]any, 0, len(pairs)+len(whereArgs))
	for _, item := range pairs {
		setClauses = append(setClauses, item.column+" = ?")
		args = append(args, item.value)
	}
	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", b.table, strings.Join(setClauses, ", "), where)
	return query, args, nil
}
