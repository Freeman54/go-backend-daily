package transactionalinbox

import (
	"errors"
	"testing"
	"time"
)

func TestClaimAndComplete(t *testing.T) {
	store := NewStore()
	store.Add("evt-1", "payload")

	now := time.Unix(100, 0)
	msg, err := store.Claim("evt-1", now, 5*time.Second)
	if err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	if msg.Attempts != 1 || msg.State != StateProcessing {
		t.Fatalf("unexpected claimed message: %#v", msg)
	}

	if err := store.Complete("evt-1"); err != nil {
		t.Fatalf("complete failed: %v", err)
	}
	saved, ok := store.Get("evt-1")
	if !ok || saved.State != StateDone {
		t.Fatalf("saved message = %#v, %v", saved, ok)
	}
}

func TestClaimRejectsActiveLease(t *testing.T) {
	store := NewStore()
	store.Add("evt-2", "payload")

	now := time.Unix(200, 0)
	if _, err := store.Claim("evt-2", now, 10*time.Second); err != nil {
		t.Fatalf("first claim failed: %v", err)
	}
	if _, err := store.Claim("evt-2", now.Add(3*time.Second), 10*time.Second); !errors.Is(err, ErrBusy) {
		t.Fatalf("claim error = %v want ErrBusy", err)
	}
}

func TestFailSchedulesRetry(t *testing.T) {
	store := NewStore()
	store.Add("evt-3", "payload")

	now := time.Unix(300, 0)
	if _, err := store.Claim("evt-3", now, 5*time.Second); err != nil {
		t.Fatalf("claim failed: %v", err)
	}
	retryAt := now.Add(30 * time.Second)
	if err := store.Fail("evt-3", retryAt, "downstream timeout"); err != nil {
		t.Fatalf("fail failed: %v", err)
	}
	if _, err := store.Claim("evt-3", now.Add(10*time.Second), 5*time.Second); !errors.Is(err, ErrBusy) {
		t.Fatalf("early retry error = %v want ErrBusy", err)
	}

	msg, err := store.Claim("evt-3", retryAt, 5*time.Second)
	if err != nil {
		t.Fatalf("retry claim failed: %v", err)
	}
	if msg.Attempts != 2 || msg.LastError != "downstream timeout" {
		t.Fatalf("retry message = %#v", msg)
	}
}
