package fencingtoken

import (
	"errors"
	"testing"
	"time"
)

func TestFencingTokenRejectsStaleWriter(t *testing.T) {
	now := time.Date(2026, 7, 3, 10, 0, 0, 0, time.UTC)
	manager := NewManager()

	first, err := manager.Grant("inventory", "node-a", now, time.Minute)
	if err != nil {
		t.Fatalf("grant first lease: %v", err)
	}

	second, err := manager.Grant("inventory", "node-b", now.Add(2*time.Minute), time.Minute)
	if err != nil {
		t.Fatalf("grant second lease: %v", err)
	}

	if err := manager.ValidateWrite("inventory", first.Token, now.Add(2*time.Minute)); !errors.Is(err, ErrStaleToken) {
		t.Fatalf("expected stale token, got %v", err)
	}

	if err := manager.ValidateWrite("inventory", second.Token, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("expected latest token to pass, got %v", err)
	}
}

func TestFencingTokenRejectsConcurrentHolder(t *testing.T) {
	now := time.Date(2026, 7, 3, 10, 0, 0, 0, time.UTC)
	manager := NewManager()

	if _, err := manager.Grant("payment", "node-a", now, time.Minute); err != nil {
		t.Fatalf("grant first lease: %v", err)
	}

	if _, err := manager.Grant("payment", "node-b", now.Add(10*time.Second), time.Minute); !errors.Is(err, ErrLeaseHeld) {
		t.Fatalf("expected lease held error, got %v", err)
	}
}

func TestFencingTokenExpiresLease(t *testing.T) {
	now := time.Date(2026, 7, 3, 10, 0, 0, 0, time.UTC)
	manager := NewManager()

	lease, err := manager.Grant("email", "node-a", now, time.Minute)
	if err != nil {
		t.Fatalf("grant lease: %v", err)
	}

	if err := manager.ValidateWrite("email", lease.Token, now.Add(2*time.Minute)); !errors.Is(err, ErrLeaseExpired) {
		t.Fatalf("expected expired lease, got %v", err)
	}
}
