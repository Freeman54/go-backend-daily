package cacherefreshjitter

import (
	"testing"
	"time"
)

func TestRefreshAfterUsesStableBoundedJitter(t *testing.T) {
	base := 10 * time.Minute
	a := RefreshAfter("tenant-a", base, 20)
	if a != RefreshAfter("tenant-a", base, 20) {
		t.Fatal("same key must be stable")
	}
	if a < 8*time.Minute || a > 12*time.Minute {
		t.Fatalf("delay %v outside bounds", a)
	}
	if a == RefreshAfter("tenant-b", base, 20) {
		t.Fatal("different keys should normally spread")
	}
}

func TestRefreshAfterClampsPercent(t *testing.T) {
	base := time.Minute
	if got := RefreshAfter("key", base, -1); got != base {
		t.Fatalf("got %v, want %v", got, base)
	}
	got := RefreshAfter("key", base, 200)
	if got < 0 || got > 2*base {
		t.Fatalf("clamped delay out of range: %v", got)
	}
}
