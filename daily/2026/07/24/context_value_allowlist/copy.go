package contextvalueallowlist

import "context"

// Copy 仅把明确允许的键从 source 复制到 destination。
func Copy(destination, source context.Context, allowed []any) context.Context {
	result := destination
	for _, key := range allowed {
		if value := source.Value(key); value != nil {
			result = context.WithValue(result, key, value)
		}
	}
	return result
}
