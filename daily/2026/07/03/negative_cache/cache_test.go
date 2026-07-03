package negativecache

import (
	"testing"
	"time"
)

func TestNegativeCachePositiveAndNegativeTTL(t *testing.T) {
	now := time.Date(2026, 7, 3, 10, 0, 0, 0, time.UTC)
	cache := New(10*time.Minute, time.Minute)
	cache.SetNow(func() time.Time { return now })

	cache.SetFound("user:1", "alice")
	cache.SetNotFound("user:2")

	if value, state := cache.Get("user:1"); value != "alice" || state != StateHit {
		t.Fatalf("expected positive hit, got value=%q state=%v", value, state)
	}

	if value, state := cache.Get("user:2"); value != "" || state != StateNegativeHit {
		t.Fatalf("expected negative hit, got value=%q state=%v", value, state)
	}

	now = now.Add(2 * time.Minute)
	if value, state := cache.Get("user:1"); value != "alice" || state != StateHit {
		t.Fatalf("expected positive entry to remain, got value=%q state=%v", value, state)
	}

	if value, state := cache.Get("user:2"); value != "" || state != StateMiss {
		t.Fatalf("expected negative entry to expire, got value=%q state=%v", value, state)
	}
}

func TestNegativeCacheCanPromoteNotFoundToFound(t *testing.T) {
	now := time.Date(2026, 7, 3, 10, 0, 0, 0, time.UTC)
	cache := New(10*time.Minute, time.Minute)
	cache.SetNow(func() time.Time { return now })

	cache.SetNotFound("order:42")
	cache.SetFound("order:42", "ready")

	if value, state := cache.Get("order:42"); value != "ready" || state != StateHit {
		t.Fatalf("expected positive entry after overwrite, got value=%q state=%v", value, state)
	}
}
