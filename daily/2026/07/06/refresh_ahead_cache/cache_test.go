package refreshaheadcache

import (
	"testing"
	"time"
)

func TestGetMissRequestsRefresh(t *testing.T) {
	cache := New(5 * time.Second)

	result := cache.Get(time.Now())
	if result.Hit {
		t.Fatal("Hit = true, want false")
	}
	if !result.Refresh {
		t.Fatal("Refresh = false, want true")
	}
}

func TestGetHotEntryDoesNotRefresh(t *testing.T) {
	now := time.Now()
	cache := New(5 * time.Second)
	cache.Store(Item{
		Value:     "profile",
		ExpiresAt: now.Add(30 * time.Second),
	})

	result := cache.Get(now)
	if !result.Hit || result.Refresh {
		t.Fatalf("result = %+v", result)
	}
}

func TestGetNearExpiryTriggersOneRefresh(t *testing.T) {
	now := time.Now()
	cache := New(5 * time.Second)
	cache.Store(Item{
		Value:     "profile",
		ExpiresAt: now.Add(3 * time.Second),
	})

	first := cache.Get(now)
	if !first.Hit || !first.Refresh {
		t.Fatalf("first = %+v", first)
	}

	second := cache.Get(now.Add(500 * time.Millisecond))
	if !second.Hit || second.Refresh {
		t.Fatalf("second = %+v", second)
	}
	if !second.Stale {
		t.Fatalf("second.Stale = false, want true")
	}
}
