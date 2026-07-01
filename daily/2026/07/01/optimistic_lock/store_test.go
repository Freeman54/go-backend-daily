package optimisticlock

import (
	"errors"
	"testing"
)

func TestUpdateBumpsVersion(t *testing.T) {
	t.Parallel()

	store := NewStore()
	created := store.Create("order-1", "pending")

	updated, err := store.Update("order-1", created.Version, "paid")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Version != 2 {
		t.Fatalf("updated version = %d, want 2", updated.Version)
	}
	if updated.Value != "paid" {
		t.Fatalf("updated value = %q, want paid", updated.Value)
	}
}

func TestUpdateRejectsStaleVersion(t *testing.T) {
	t.Parallel()

	store := NewStore()
	store.Create("order-1", "pending")

	_, err := store.Update("order-1", 0, "paid")
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("Update() error = %v, want ErrConflict", err)
	}
}

func TestGetReturnsNotFound(t *testing.T) {
	t.Parallel()

	store := NewStore()
	_, err := store.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}
