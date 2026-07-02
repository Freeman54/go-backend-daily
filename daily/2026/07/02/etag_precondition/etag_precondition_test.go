package etagprecondition

import (
	"errors"
	"testing"
)

func TestUpdateMovesVersionForward(t *testing.T) {
	t.Parallel()

	store := NewStore("draft")
	before := store.Snapshot()
	updated, err := store.Update(ETag(before.Version), "published")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Body != "published" {
		t.Fatalf("Body = %q, want published", updated.Body)
	}
	if updated.Version != before.Version+1 {
		t.Fatalf("Version = %d, want %d", updated.Version, before.Version+1)
	}
}

func TestUpdateRejectsStaleETag(t *testing.T) {
	t.Parallel()

	store := NewStore("draft")
	snapshot := store.Snapshot()
	if _, err := store.Update(ETag(snapshot.Version), "published"); err != nil {
		t.Fatalf("first Update() error = %v", err)
	}

	_, err := store.Update(ETag(snapshot.Version), "overwrite")
	if !errors.Is(err, ErrPreconditionFailed) {
		t.Fatalf("second Update() error = %v, want ErrPreconditionFailed", err)
	}
}

func TestUpdateRejectsMalformedETag(t *testing.T) {
	t.Parallel()

	store := NewStore("draft")
	_, err := store.Update("version-1", "published")
	if !errors.Is(err, ErrMalformedETag) {
		t.Fatalf("Update() error = %v, want ErrMalformedETag", err)
	}
}
