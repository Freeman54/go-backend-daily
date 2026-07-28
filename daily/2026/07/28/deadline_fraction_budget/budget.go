package deadline_fraction_budget

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoDeadline      = errors.New("parent has no deadline")
	ErrInvalidFraction = errors.New("fraction must be in (0, 1]")
)

// WithFraction 把父上下文剩余时间的一部分分配给当前阶段。
func WithFraction(parent context.Context, now time.Time, fraction float64, minimum time.Duration) (context.Context, context.CancelFunc, error) {
	deadline, ok := parent.Deadline()
	if !ok {
		return nil, nil, ErrNoDeadline
	}
	if fraction <= 0 || fraction > 1 || minimum < 0 {
		return nil, nil, ErrInvalidFraction
	}
	remaining := deadline.Sub(now)
	budget := time.Duration(float64(remaining) * fraction)
	if budget < minimum {
		budget = minimum
	}
	if budget > remaining {
		budget = remaining
	}
	ctx, cancel := context.WithDeadline(parent, now.Add(budget))
	return ctx, cancel, nil
}
