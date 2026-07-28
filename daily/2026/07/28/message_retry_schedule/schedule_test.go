package message_retry_schedule

import (
	"errors"
	"testing"
	"time"
)

func TestDelay_UsesExponentialBackoffWithDeterministicJitter(t *testing.T) {
	policy := Policy{Base: time.Second, Max: 10 * time.Second, Jitter: 0.2}
	got, err := policy.Delay(3, 0.75)
	if err != nil {
		t.Fatal(err)
	}
	if got != 4400*time.Millisecond {
		t.Fatalf("Delay() = %v, want 4.4s", got)
	}
}

func TestDelay_CapsBeforeApplyingJitter(t *testing.T) {
	policy := Policy{Base: 2 * time.Second, Max: 5 * time.Second, Jitter: 0}
	got, err := policy.Delay(10, 0.5)
	if err != nil || got != 5*time.Second {
		t.Fatalf("Delay() = (%v, %v), want (5s, nil)", got, err)
	}
}

func TestDelay_RejectsInvalidPolicyAndAttempt(t *testing.T) {
	if _, err := (Policy{}).Delay(1, 0.5); !errors.Is(err, ErrInvalidPolicy) {
		t.Fatalf("error = %v, want ErrInvalidPolicy", err)
	}
	policy := Policy{Base: time.Second, Max: time.Second}
	if _, err := policy.Delay(0, 0.5); !errors.Is(err, ErrInvalidAttempt) {
		t.Fatalf("error = %v, want ErrInvalidAttempt", err)
	}
}
