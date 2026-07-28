package sql_placeholder_rebinder

import (
	"errors"
	"strconv"
	"strings"
)

var ErrUnterminatedQuote = errors.New("unterminated SQL string")

// Rebind 把未处于单引号字符串中的 ? 转换为 PostgreSQL 的 $n 占位符。
func Rebind(query string) (string, int, error) {
	var out strings.Builder
	out.Grow(len(query))
	inQuote := false
	count := 0
	for i := 0; i < len(query); i++ {
		ch := query[i]
		if ch == '\'' {
			out.WriteByte(ch)
			if inQuote && i+1 < len(query) && query[i+1] == '\'' {
				i++
				out.WriteByte(query[i])
				continue
			}
			inQuote = !inQuote
			continue
		}
		if ch == '?' && !inQuote {
			count++
			out.WriteByte('$')
			out.WriteString(strconv.Itoa(count))
			continue
		}
		out.WriteByte(ch)
	}
	if inQuote {
		return "", 0, ErrUnterminatedQuote
	}
	return out.String(), count, nil
}
