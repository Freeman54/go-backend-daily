package ttljittercache

import (
	"math/rand"
	"testing"
	"time"
)

func TestJitterTTLWithinBounds(t *testing.T) {
	t.Parallel()

	base := 5 * time.Minute
	jitter := 2 * time.Minute
	got := JitterTTL(base, jitter, rand.New(rand.NewSource(7)))
	if got < base || got > base+jitter {
		t.Fatalf("JitterTTL() = %v, want in [%v, %v]", got, base, base+jitter)
	}
}

func TestStaggerExpirationsSpreadsKeys(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	got := StaggerExpirations(
		now,
		[]string{"a", "b", "c"},
		10*time.Minute,
		5*time.Minute,
		rand.New(rand.NewSource(1)),
	)
	if len(got) != 3 {
		t.Fatalf("len(StaggerExpirations()) = %d, want 3", len(got))
	}
	if got["a"] == got["b"] && got["b"] == got["c"] {
		t.Fatalf("all expirations are identical, want spread")
	}
	for key, expiresAt := range got {
		if expiresAt.Before(now.Add(10*time.Minute)) || expiresAt.After(now.Add(15*time.Minute)) {
			t.Fatalf("%s expiration = %v, out of bounds", key, expiresAt)
		}
	}
}

func TestJitterTTLWithoutRandomnessUsesBase(t *testing.T) {
	t.Parallel()

	base := 30 * time.Second
	if got := JitterTTL(base, time.Minute, nil); got != base {
		t.Fatalf("JitterTTL() = %v, want %v", got, base)
	}
}
