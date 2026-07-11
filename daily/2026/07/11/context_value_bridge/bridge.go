package contextvaluebridge

import (
	"context"
	"time"
)

func Detach(parent context.Context, timeout time.Duration, keys ...any) (context.Context, context.CancelFunc) {
	base := context.Background()
	for _, key := range keys {
		if key == nil {
			continue
		}
		if value := parent.Value(key); value != nil {
			base = context.WithValue(base, key, value)
		}
	}

	if timeout > 0 {
		return context.WithTimeout(base, timeout)
	}
	return context.WithCancel(base)
}
