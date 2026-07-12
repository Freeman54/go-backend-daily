package cachewritethrough

import (
	"errors"
	"testing"
)

func TestSaveWritesRepoAndCache(t *testing.T) {
	repo := NewMemoryStore()
	cache := NewMemoryStore()
	store := New(repo, cache)

	if err := store.Save("user:1", "alice", 2); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if record, ok := repo.Get("user:1"); !ok || record.Value != "alice" || record.Version != 2 {
		t.Fatalf("repo record = %#v, %v", record, ok)
	}
	if record, ok := cache.Get("user:1"); !ok || record.Value != "alice" || record.Version != 2 {
		t.Fatalf("cache record = %#v, %v", record, ok)
	}
}

func TestSaveRejectsStaleVersion(t *testing.T) {
	repo := NewMemoryStore()
	cache := NewMemoryStore()
	store := New(repo, cache)

	if err := store.Save("profile", "v2", 2); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}
	err := store.Save("profile", "v1", 1)
	if !errors.Is(err, ErrStaleWrite) {
		t.Fatalf("save error = %v want ErrStaleWrite", err)
	}
}

func TestLoadBackfillsCache(t *testing.T) {
	repo := NewMemoryStore()
	cache := NewMemoryStore()
	store := New(repo, cache)

	repo.Upsert("settings", Record{Value: "enabled", Version: 3})
	cache.Delete("settings")

	record, ok := store.Load("settings")
	if !ok || record.Value != "enabled" {
		t.Fatalf("load = %#v, %v", record, ok)
	}

	cached, ok := cache.Get("settings")
	if !ok || cached.Version != 3 {
		t.Fatalf("cache backfill = %#v, %v", cached, ok)
	}
}
