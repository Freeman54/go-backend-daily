package txdeadlockretry

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stateErr struct {
	state string
	msg   string
}

func (e stateErr) Error() string    { return e.msg }
func (e stateErr) SQLState() string { return e.state }

func TestPolicyRetriesDeadlockThenSucceeds(t *testing.T) {
	policy := Policy{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 2 * time.Millisecond}
	attempts := 0
	var sleeps []time.Duration

	err := policy.Run(context.Background(), func(context.Context) error {
		attempts++
		if attempts < 3 {
			return stateErr{state: "40P01", msg: "deadlock detected"}
		}
		return nil
	}, func(d time.Duration) {
		sleeps = append(sleeps, d)
	})
	if err != nil {
		t.Fatalf("expected eventual success, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if len(sleeps) != 2 || sleeps[0] != time.Millisecond || sleeps[1] != 2*time.Millisecond {
		t.Fatalf("unexpected backoff sequence: %v", sleeps)
	}
}

func TestPolicyStopsOnNonRetryableError(t *testing.T) {
	policy := Policy{MaxAttempts: 4, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	attempts := 0
	expected := errors.New("syntax error")

	err := policy.Run(context.Background(), func(context.Context) error {
		attempts++
		return expected
	}, nil)
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	if attempts != 1 {
		t.Fatalf("expected one attempt, got %d", attempts)
	}
}

func TestPolicyHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	policy := Policy{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
	err := policy.Run(ctx, func(context.Context) error {
		t.Fatal("fn should not be called after cancellation")
		return nil
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}
