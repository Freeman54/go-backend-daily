package writebehindcache

import "testing"

func TestFlushAppliesQueuedUpdates(t *testing.T) {
	store := NewStore()
	store.Enqueue("user:1", "v1")
	store.Enqueue("user:2", "v2")

	batch := store.Flush(10)

	if len(batch) != 2 {
		t.Fatalf("unexpected batch size: %d", len(batch))
	}
	if got, ok := store.Get("user:1"); !ok || got != "v1" {
		t.Fatalf("unexpected value after flush: %q %v", got, ok)
	}
}

func TestFlushCanBePartial(t *testing.T) {
	store := NewStore()
	store.Enqueue("user:1", "v1")
	store.Enqueue("user:2", "v2")

	batch := store.Flush(1)

	if len(batch) != 1 {
		t.Fatalf("unexpected batch size: %d", len(batch))
	}
	if pending := store.Pending(); pending != 1 {
		t.Fatalf("expected one pending item, got %d", pending)
	}
	if _, ok := store.Get("user:2"); ok {
		t.Fatal("second item should not be visible before next flush")
	}
}

func TestLaterUpdateWins(t *testing.T) {
	store := NewStore()
	store.Enqueue("user:1", "v1")
	store.Enqueue("user:1", "v2")

	store.Flush(10)

	if got, _ := store.Get("user:1"); got != "v2" {
		t.Fatalf("expected latest value, got %q", got)
	}
}
