package contextbudgetreserve

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNoDeadline         = errors.New("context has no deadline")
	ErrInsufficientBudget = errors.New("insufficient context budget")
)

func WithReserve(parent context.Context, reserve time.Duration) (context.Context, context.CancelFunc, error) {
	if reserve < 0 {
		return nil, nil, fmt.Errorf("reserve must not be negative")
	}
	deadline, ok := parent.Deadline()
	if !ok {
		return nil, nil, ErrNoDeadline
	}
	childDeadline := deadline.Add(-reserve)
	if !childDeadline.After(time.Now()) {
		return nil, nil, ErrInsufficientBudget
	}
	ctx, cancel := context.WithDeadline(parent, childDeadline)
	return ctx, cancel, nil
}
