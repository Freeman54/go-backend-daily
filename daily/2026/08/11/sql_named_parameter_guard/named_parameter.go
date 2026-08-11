// Package namedparameter validates named SQL argument contracts.
package namedparameter

import (
	"errors"
	"sort"
	"unicode"
)

var ErrArgumentMismatch = errors.New("SQL named argument mismatch")

// Validate checks that every :name placeholder has exactly one supplied value
// and that no unused values are supplied. Colons inside quoted strings are ignored.
func Validate(query string, args map[string]any) error {
	wanted := make(map[string]struct{})
	for i, quote := 0, rune(0); i < len(query); {
		r := rune(query[i])
		if quote != 0 {
			if r == quote {
				quote = 0
			}
			i++
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			i++
			continue
		}
		if r != ':' || i+1 >= len(query) || !isStart(rune(query[i+1])) {
			i++
			continue
		}
		j := i + 2
		for j < len(query) && isPart(rune(query[j])) {
			j++
		}
		wanted[query[i+1:j]] = struct{}{}
		i = j
	}
	if !sameKeys(wanted, args) {
		return ErrArgumentMismatch
	}
	return nil
}

func Names(args map[string]any) []string {
	names := make([]string, 0, len(args))
	for name := range args {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func isStart(r rune) bool { return r == '_' || unicode.IsLetter(r) }
func isPart(r rune) bool  { return isStart(r) || unicode.IsDigit(r) }

func sameKeys(wanted map[string]struct{}, args map[string]any) bool {
	if len(wanted) != len(args) {
		return false
	}
	for name := range wanted {
		if _, ok := args[name]; !ok {
			return false
		}
	}
	return true
}
