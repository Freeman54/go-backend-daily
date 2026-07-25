package deadline_cause_policy

import (
	"context"
	"errors"
	"time"
)

func WithTimeoutCause(parent context.Context, timeout time.Duration, cause error) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(parent, timeout, cause)
}

func NewTimeoutCause(parent context.Context, timeout time.Duration, cause error) (context.Context, context.CancelFunc, error) {
	if timeout <= 0 {
		return nil, nil, errors.New("timeout must be positive")
	}
	if cause == nil {
		return nil, nil, errors.New("cause must not be nil")
	}
	ctx, cancel := WithTimeoutCause(parent, timeout, cause)
	return ctx, cancel, nil
}
