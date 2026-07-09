package versionedcachekeys

import "testing"

func TestStoreReadsCurrentVersion(t *testing.T) {
	store := New()
	key := store.Put("catalog", "42", "book")

	if key != "catalog:v0:42" {
		t.Fatalf("unexpected key %q", key)
	}
	got, ok := store.Get("catalog", "42")
	if !ok || got != "book" {
		t.Fatalf("expected current value, got %q, ok=%v", got, ok)
	}
}

func TestBumpInvalidatesOldKeysWithoutDeleteSweep(t *testing.T) {
	store := New()
	store.Put("catalog", "42", "old")

	version := store.Bump("catalog")
	if version != 1 {
		t.Fatalf("expected version 1, got %d", version)
	}
	if _, ok := store.Get("catalog", "42"); ok {
		t.Fatal("expected old key to be invisible after version bump")
	}

	key := store.Put("catalog", "42", "new")
	if key != "catalog:v1:42" {
		t.Fatalf("unexpected bumped key %q", key)
	}
	got, ok := store.Get("catalog", "42")
	if !ok || got != "new" {
		t.Fatalf("expected new value, got %q, ok=%v", got, ok)
	}
}

func TestNamespacesAreIsolated(t *testing.T) {
	store := New()
	store.Put("catalog", "42", "book")
	store.Put("profile", "42", "alice")
	store.Bump("catalog")

	got, ok := store.Get("profile", "42")
	if !ok || got != "alice" {
		t.Fatalf("expected other namespace to stay readable, got %q, ok=%v", got, ok)
	}
}
