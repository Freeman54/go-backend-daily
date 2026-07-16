package cachetombstone

import (
	"testing"
	"time"
)

func TestDeleteBlocksOlderPut(t *testing.T) {
	now := time.Date(2026, 7, 16, 0, 0, 0, 0, time.UTC)
	store := New(func() time.Time { return now })
	if !store.Put("user:1", "old", 1) {
		t.Fatal("expected initial put to succeed")
	}

	store.Delete("user:1", 2, time.Minute)

	if store.Put("user:1", "stale", 1) {
		t.Fatal("expected stale put to be rejected")
	}
}

func TestHigherVersionCanRecoverFromTombstone(t *testing.T) {
	now := time.Date(2026, 7, 16, 0, 0, 0, 0, time.UTC)
	store := New(func() time.Time { return now })
	store.Delete("user:1", 2, time.Minute)

	if !store.Put("user:1", "fresh", 3) {
		t.Fatal("expected higher version put to succeed")
	}

	got, ok := store.Get("user:1")
	if !ok || got != "fresh" {
		t.Fatalf("unexpected get result: ok=%v value=%q", ok, got)
	}
}

func TestExpiredTombstoneAllowsRewrite(t *testing.T) {
	now := time.Date(2026, 7, 16, 0, 0, 0, 0, time.UTC)
	store := New(func() time.Time { return now })
	store.Delete("user:1", 2, time.Minute)

	now = now.Add(2 * time.Minute)

	if !store.Put("user:1", "reloaded", 1) {
		t.Fatal("expected put after expiry to succeed")
	}
}
