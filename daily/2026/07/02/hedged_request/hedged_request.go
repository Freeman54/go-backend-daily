package hedgedrequest

import (
	"context"
	"errors"
	"time"
)

var ErrNoSuccess = errors.New("no backend succeeded")

type Backend func(context.Context) (string, error)

type result struct {
	value string
	err   error
}

func Do(ctx context.Context, primary Backend, hedge Backend, delay time.Duration) (string, error) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan result, 2)
	go func() {
		value, err := primary(childCtx)
		results <- result{value: value, err: err}
	}()

	timer := time.NewTimer(delay)
	defer timer.Stop()

	hedgeStarted := false
	failures := make([]error, 0, 2)
	for len(failures) < 2 {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timer.C:
			if hedgeStarted {
				continue
			}
			hedgeStarted = true
			go func() {
				value, err := hedge(childCtx)
				results <- result{value: value, err: err}
			}()
		case outcome := <-results:
			if outcome.err == nil {
				cancel()
				return outcome.value, nil
			}
			failures = append(failures, outcome.err)
			if !hedgeStarted {
				hedgeStarted = true
				go func() {
					value, err := hedge(childCtx)
					results <- result{value: value, err: err}
				}()
			}
		}
	}

	return "", errors.Join(append([]error{ErrNoSuccess}, failures...)...)
}
