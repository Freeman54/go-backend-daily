package tokenbucket

import (
	"testing"
	"time"
)

func TestAllowNConsumesAvailableTokens(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	bucket := New(5, time.Second, now)

	if !bucket.AllowN(now, 3) {
		t.Fatalf("AllowN() = false, want true")
	}
	if bucket.AllowN(now, 3) {
		t.Fatalf("AllowN() = true, want false")
	}
}

func TestAllowNRefillsOverTime(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	bucket := New(4, 2*time.Second, now)
	if !bucket.AllowN(now, 4) {
		t.Fatalf("initial AllowN() = false, want true")
	}

	later := now.Add(time.Second)
	if !bucket.AllowN(later, 2) {
		t.Fatalf("AllowN() after refill = false, want true")
	}
	if bucket.AllowN(later, 1) {
		t.Fatalf("AllowN() should reject when token is exhausted")
	}
}

func TestTokensCapsAtCapacity(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	bucket := New(3, time.Second, now)
	if !bucket.AllowN(now, 3) {
		t.Fatalf("AllowN() = false, want true")
	}

	got := bucket.Tokens(now.Add(5 * time.Second))
	if got != 3 {
		t.Fatalf("Tokens() = %v, want 3", got)
	}
}
