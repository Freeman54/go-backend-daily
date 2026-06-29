package retrybackoff

import (
	"context"
	"errors"
	"time"
)

type RetryableError struct {
	Err error
}

func (e RetryableError) Error() string {
	return e.Err.Error()
}

func (e RetryableError) Unwrap() error {
	return e.Err
}

type Policy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func (p Policy) Delay(attempt int) time.Duration {
	if attempt <= 1 {
		return 0
	}

	delay := p.BaseDelay << (attempt - 2)
	if p.MaxDelay > 0 && delay > p.MaxDelay {
		return p.MaxDelay
	}
	return delay
}

func Do(ctx context.Context, policy Policy, sleep func(context.Context, time.Duration) error, op func(context.Context) error) error {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 1
	}
	if sleep == nil {
		sleep = defaultSleep
	}

	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if attempt > 1 {
			if err := sleep(ctx, policy.Delay(attempt)); err != nil {
				return err
			}
		}

		err := op(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		var retryable RetryableError
		if !errors.As(err, &retryable) {
			return err
		}
	}

	return lastErr
}

func defaultSleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
