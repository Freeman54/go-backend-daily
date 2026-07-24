package sqllockorder

import (
	"fmt"
	"strings"
)

// Plan 返回全局一致的加锁顺序，减少事务间循环等待。
func Plan(ids []int64) ([]int64, error) {
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return nil, fmt.Errorf("resource ID must be positive: %d", id)
		}
		seen[id] = struct{}{}
	}
	result := make([]int64, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j] < result[j-1]; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result, nil
}

func ForUpdateQuery(table string, ids []int64) (string, []int64, error) {
	if !validIdentifier(table) {
		return "", nil, fmt.Errorf("invalid table name %q", table)
	}
	ordered, err := Plan(ids)
	if err != nil {
		return "", nil, err
	}
	if len(ordered) == 0 {
		return "", nil, fmt.Errorf("at least one resource ID is required")
	}
	placeholders := make([]string, len(ordered))
	for i := range ordered {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("SELECT id FROM %s WHERE id IN (%s) ORDER BY id FOR UPDATE", table, strings.Join(placeholders, ","))
	return query, ordered, nil
}

func validIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for i, r := range value {
		if !(r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || i > 0 && r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}
