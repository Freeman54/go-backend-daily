package txretry

import (
	"context"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Runner struct {
	MaxAttempts int
	Begin       func(context.Context) (Tx, error)
	IsRetryable func(error) bool
}

func (r Runner) Do(ctx context.Context, fn func(context.Context, Tx) error) error {
	maxAttempts := r.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tx, err := r.Begin(ctx)
		if err != nil {
			return err
		}

		err = fn(ctx, tx)
		if err != nil {
			_ = tx.Rollback()
			lastErr = err
			if r.IsRetryable != nil && r.IsRetryable(err) && attempt < maxAttempts {
				continue
			}
			return err
		}

		if err = tx.Commit(); err != nil {
			_ = tx.Rollback()
			lastErr = err
			if r.IsRetryable != nil && r.IsRetryable(err) && attempt < maxAttempts {
				continue
			}
			return err
		}
		return nil
	}

	return lastErr
}
