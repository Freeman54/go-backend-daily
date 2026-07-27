package message_expiry_policy

import (
	"errors"
	"testing"
	"time"
)

func TestPolicy_RemainingTTLAccountsForTimeAlreadySpent(t *testing.T) {
	policy, err := New(10 * time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	createdAt := time.Unix(100, 0)
	ttl, err := policy.RemainingTTL(createdAt, createdAt.Add(3*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if ttl != 7*time.Minute {
		t.Fatalf("RemainingTTL() = %v, want 7m", ttl)
	}
}

func TestPolicy_RejectsExpiredMessage(t *testing.T) {
	policy, _ := New(5 * time.Minute)
	createdAt := time.Unix(200, 0)
	_, err := policy.RemainingTTL(createdAt, createdAt.Add(5*time.Minute))
	if !errors.Is(err, ErrExpired) {
		t.Fatalf("RemainingTTL() error = %v, want ErrExpired", err)
	}
}

func TestPolicy_RejectsClockBeforeCreation(t *testing.T) {
	policy, _ := New(time.Minute)
	createdAt := time.Unix(300, 0)
	_, err := policy.RemainingTTL(createdAt, createdAt.Add(-time.Second))
	if !errors.Is(err, ErrInvalidTimestamp) {
		t.Fatalf("RemainingTTL() error = %v, want ErrInvalidTimestamp", err)
	}
}

func TestNew_RejectsNonPositiveTTL(t *testing.T) {
	if _, err := New(0); err == nil {
		t.Fatal("zero TTL should fail")
	}
}
