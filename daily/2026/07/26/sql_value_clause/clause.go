package sql_value_clause

import (
	"errors"
	"strconv"
	"strings"
)

// Build 创建 PostgreSQL 批量 INSERT 的 VALUES 占位符片段。
func Build(columns, rows, start int) (clause string, next int, err error) {
	if columns <= 0 || rows <= 0 || start <= 0 {
		return "", start, errors.New("columns, rows and start must be positive")
	}
	var builder strings.Builder
	placeholder := start
	for row := 0; row < rows; row++ {
		if row > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('(')
		for column := 0; column < columns; column++ {
			if column > 0 {
				builder.WriteByte(',')
			}
			builder.WriteByte('$')
			builder.WriteString(strconv.Itoa(placeholder))
			placeholder++
		}
		builder.WriteByte(')')
	}
	return builder.String(), placeholder, nil
}
