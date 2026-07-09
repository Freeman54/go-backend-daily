package txdeadlockretry

import (
	"context"
	"errors"
	"strings"
	"time"
)

type sqlStateCarrier interface {
	SQLState() string
}

type Policy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func (p Policy) Run(ctx context.Context, fn func(context.Context) error, sleep func(time.Duration)) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = 10 * time.Millisecond
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 200 * time.Millisecond
	}
	if sleep == nil {
		sleep = time.Sleep
	}

	var lastErr error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}
		if !ShouldRetry(lastErr) || attempt == p.MaxAttempts {
			return lastErr
		}

		sleep(backoff(attempt, p.BaseDelay, p.MaxDelay))
	}

	return lastErr
}

func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	var carrier sqlStateCarrier
	if errors.As(err, &carrier) {
		switch carrier.SQLState() {
		case "40001", "40P01":
			return true
		}
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "deadlock") || strings.Contains(message, "lock wait timeout")
}

func backoff(attempt int, base, max time.Duration) time.Duration {
	delay := base << (attempt - 1)
	if delay > max {
		return max
	}
	return delay
}
