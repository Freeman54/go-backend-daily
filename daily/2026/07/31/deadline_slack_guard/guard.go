package deadlineslackguard

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoDeadline        = errors.New("context has no deadline")
	ErrInsufficientSlack = errors.New("deadline slack is insufficient")
	ErrInvalidMinimum    = errors.New("minimum slack must be positive")
)

// Check 判断上下文剩余时间能否安全启动下一阶段工作。
func Check(ctx context.Context, now time.Time, minimum time.Duration) error {
	if minimum <= 0 {
		return ErrInvalidMinimum
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrNoDeadline
	}
	if deadline.Sub(now) < minimum {
		return ErrInsufficientSlack
	}
	return nil
}
