package cache_probabilistic_refresh

import (
	"testing"
	"time"
)

func TestShouldRefresh_RefreshesExpiredEntry(t *testing.T) {
	now := time.Unix(100, 0)
	if !ShouldRefresh(now, now, 10*time.Second, 0.9, 0.2) {
		t.Fatal("expired entry should refresh")
	}
}

func TestShouldRefresh_UsesRandomizedEarlyWindow(t *testing.T) {
	now := time.Unix(100, 0)
	expires := now.Add(2 * time.Second)
	if !ShouldRefresh(now, expires, 10*time.Second, 0.9, 0.2) {
		t.Fatal("high sample should refresh inside early window")
	}
	if ShouldRefresh(now, expires, 10*time.Second, 0.1, 0.2) {
		t.Fatal("low sample should not refresh yet")
	}
}

func TestShouldRefresh_RejectsInvalidInputs(t *testing.T) {
	now := time.Now()
	if ShouldRefresh(now, now.Add(time.Second), 0, 0.5, 0.2) {
		t.Fatal("non-positive ttl must not trigger early refresh")
	}
	if ShouldRefresh(now, now.Add(time.Second), time.Second, 1.1, 0.2) {
		t.Fatal("sample outside [0,1] must not refresh")
	}
}
