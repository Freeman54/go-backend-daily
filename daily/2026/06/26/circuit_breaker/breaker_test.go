package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestBreakerOpensAfterThresholdAndHalfOpensAfterCooldown(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	breaker := New(Config{
		FailureThreshold: 2,
		HalfOpenAfter:    time.Minute,
		Now:              func() time.Time { return now },
	})

	failure := errors.New("upstream down")
	if err := breaker.Execute(func() error { return failure }); !errors.Is(err, failure) {
		t.Fatalf("first Execute() error = %v, want %v", err, failure)
	}
	if err := breaker.Execute(func() error { return failure }); !errors.Is(err, failure) {
		t.Fatalf("second Execute() error = %v, want %v", err, failure)
	}

	if state := breaker.State(); state != StateOpen {
		t.Fatalf("State() = %v, want %v", state, StateOpen)
	}

	err := breaker.Execute(func() error { return nil })
	if !errors.Is(err, ErrOpen) {
		t.Fatalf("Execute() while open = %v, want ErrOpen", err)
	}

	now = now.Add(2 * time.Minute)

	if state := breaker.State(); state != StateHalfOpen {
		t.Fatalf("State() after cooldown = %v, want %v", state, StateHalfOpen)
	}

	if err := breaker.Execute(func() error { return nil }); err != nil {
		t.Fatalf("half-open Execute() returned error: %v", err)
	}

	if state := breaker.State(); state != StateClosed {
		t.Fatalf("State() after success = %v, want %v", state, StateClosed)
	}
}
