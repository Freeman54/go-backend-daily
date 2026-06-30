package stalewhilerevalidate

import (
	"errors"
	"testing"
	"time"
)

func TestGetReturnsFreshValue(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	cache := New()
	cache.Set("profile:1", Entry{
		Value:      "fresh",
		ExpiresAt:  now.Add(time.Minute),
		StaleUntil: now.Add(2 * time.Minute),
	})

	got, err := cache.Get(now, "profile:1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !got.Fresh || got.ShouldRefresh {
		t.Fatalf("Get() = %+v, want fresh result", got)
	}
}

func TestGetReturnsStaleAndStartsSingleRefresh(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	cache := New()
	cache.Set("profile:1", Entry{
		Value:      "stale",
		ExpiresAt:  now.Add(-time.Second),
		StaleUntil: now.Add(time.Minute),
	})

	first, err := cache.Get(now, "profile:1")
	if err != nil {
		t.Fatalf("first Get() error = %v", err)
	}
	second, err := cache.Get(now, "profile:1")
	if err != nil {
		t.Fatalf("second Get() error = %v", err)
	}

	if !first.ShouldRefresh {
		t.Fatalf("first Get() should request refresh")
	}
	if second.ShouldRefresh {
		t.Fatalf("second Get() should not duplicate refresh")
	}
}

func TestRefreshReplacesEntry(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	cache := New()
	cache.Set("profile:1", Entry{
		Value:      "stale",
		ExpiresAt:  now.Add(-time.Second),
		StaleUntil: now.Add(time.Minute),
	})

	value, err := cache.Refresh(now, "profile:1", time.Minute, 2*time.Minute, func() (string, error) {
		return "fresh", nil
	})
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if value != "fresh" {
		t.Fatalf("Refresh() value = %q, want fresh", value)
	}

	got, err := cache.Get(now.Add(30*time.Second), "profile:1")
	if err != nil {
		t.Fatalf("Get() after refresh error = %v", err)
	}
	if got.Value != "fresh" || !got.Fresh {
		t.Fatalf("Get() after refresh = %+v, want fresh entry", got)
	}
}

func TestRefreshClearsRefreshingFlagOnFailure(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	cache := New()
	cache.Set("profile:1", Entry{
		Value:      "stale",
		ExpiresAt:  now.Add(-time.Second),
		StaleUntil: now.Add(time.Minute),
	})

	if _, err := cache.Get(now, "profile:1"); err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	boom := errors.New("db timeout")
	_, err := cache.Refresh(now, "profile:1", time.Minute, time.Minute, func() (string, error) {
		return "", boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("Refresh() error = %v, want %v", err, boom)
	}

	got, err := cache.Get(now, "profile:1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !got.ShouldRefresh {
		t.Fatalf("Get() should request refresh again after failure")
	}
}
