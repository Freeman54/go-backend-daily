package retrybackoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoRetriesRetryableErrors(t *testing.T) {
	t.Parallel()

	var attempts int
	var delays []time.Duration

	err := Do(
		context.Background(),
		Policy{MaxAttempts: 4, BaseDelay: 10 * time.Millisecond, MaxDelay: 40 * time.Millisecond},
		func(_ context.Context, d time.Duration) error {
			delays = append(delays, d)
			return nil
		},
		func(context.Context) error {
			attempts++
			if attempts < 3 {
				return RetryableError{Err: errors.New("temporary unavailable")}
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
	want := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}
	if len(delays) != len(want) {
		t.Fatalf("delay count = %d, want %d", len(delays), len(want))
	}
	for i := range want {
		if delays[i] != want[i] {
			t.Fatalf("delay[%d] = %v, want %v", i, delays[i], want[i])
		}
	}
}

func TestDoStopsOnPermanentError(t *testing.T) {
	t.Parallel()

	var attempts int
	permanent := errors.New("bad request")

	err := Do(context.Background(), Policy{MaxAttempts: 5, BaseDelay: time.Millisecond}, nil, func(context.Context) error {
		attempts++
		return permanent
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("Do() error = %v, want %v", err, permanent)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestDoReturnsContextErrorWhileSleeping(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var attempts int
	err := Do(ctx, Policy{MaxAttempts: 3, BaseDelay: time.Second}, func(ctx context.Context, _ time.Duration) error {
		cancel()
		<-ctx.Done()
		return ctx.Err()
	}, func(context.Context) error {
		attempts++
		return RetryableError{Err: errors.New("timeout")}
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
