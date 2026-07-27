package cache_stale_if_error

import (
	"errors"
	"testing"
	"time"
)

func TestCache_GetReturnsFreshValueWithoutLoading(t *testing.T) {
	now := time.Unix(100, 0)
	cache, _ := New[string](time.Minute, 5*time.Minute)
	cache.Store("user:1", "cached", now)
	calls := 0

	got, err := cache.Get("user:1", now.Add(30*time.Second), func() (string, error) {
		calls++
		return "loaded", nil
	})
	if err != nil || got != "cached" || calls != 0 {
		t.Fatalf("Get() = (%q, %v), calls=%d", got, err, calls)
	}
}

func TestCache_GetUsesStaleValueWhenRefreshFails(t *testing.T) {
	now := time.Unix(200, 0)
	cache, _ := New[string](time.Minute, 5*time.Minute)
	cache.Store("user:1", "stale", now)
	loadErr := errors.New("upstream unavailable")

	got, err := cache.Get("user:1", now.Add(2*time.Minute), func() (string, error) {
		return "", loadErr
	})
	if err != nil || got != "stale" {
		t.Fatalf("Get() = (%q, %v), want stale value", got, err)
	}
}

func TestCache_GetReturnsLoaderErrorAfterStaleWindow(t *testing.T) {
	now := time.Unix(300, 0)
	cache, _ := New[string](time.Minute, 5*time.Minute)
	cache.Store("user:1", "expired", now)
	loadErr := errors.New("upstream unavailable")

	_, err := cache.Get("user:1", now.Add(7*time.Minute), func() (string, error) {
		return "", loadErr
	})
	if !errors.Is(err, loadErr) {
		t.Fatalf("Get() error = %v, want loader error", err)
	}
}

func TestNew_RejectsInvalidDurations(t *testing.T) {
	if _, err := New[string](0, time.Minute); err == nil {
		t.Fatal("zero fresh TTL should fail")
	}
	if _, err := New[string](time.Minute, 30*time.Second); err == nil {
		t.Fatal("stale TTL shorter than fresh TTL should fail")
	}
}
