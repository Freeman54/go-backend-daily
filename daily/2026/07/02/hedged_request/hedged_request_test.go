package hedgedrequest

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoReturnsPrimaryWithoutLaunchingHedge(t *testing.T) {
	t.Parallel()

	var hedgeCalls atomic.Int32
	got, err := Do(
		context.Background(),
		func(context.Context) (string, error) {
			return "primary", nil
		},
		func(context.Context) (string, error) {
			hedgeCalls.Add(1)
			return "hedge", nil
		},
		50*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got != "primary" {
		t.Fatalf("Do() = %q, want primary", got)
	}
	if hedgeCalls.Load() != 0 {
		t.Fatalf("hedge called %d times, want 0", hedgeCalls.Load())
	}
}

func TestDoReturnsFastHedge(t *testing.T) {
	t.Parallel()

	start := time.Now()
	got, err := Do(
		context.Background(),
		func(ctx context.Context) (string, error) {
			select {
			case <-time.After(200 * time.Millisecond):
				return "primary", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
		func(ctx context.Context) (string, error) {
			select {
			case <-time.After(20 * time.Millisecond):
				return "hedge", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
		30*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if got != "hedge" {
		t.Fatalf("Do() = %q, want hedge", got)
	}
	if elapsed := time.Since(start); elapsed >= 120*time.Millisecond {
		t.Fatalf("Do() took %v, want < 120ms", elapsed)
	}
}

func TestDoJoinsFailures(t *testing.T) {
	t.Parallel()

	primaryErr := errors.New("primary failed")
	hedgeErr := errors.New("hedge failed")
	_, err := Do(
		context.Background(),
		func(context.Context) (string, error) { return "", primaryErr },
		func(context.Context) (string, error) { return "", hedgeErr },
		time.Hour,
	)
	if !errors.Is(err, ErrNoSuccess) {
		t.Fatalf("Do() error = %v, want ErrNoSuccess", err)
	}
	if !errors.Is(err, primaryErr) || !errors.Is(err, hedgeErr) {
		t.Fatalf("Do() error = %v, want joined backend errors", err)
	}
}
