package prioritylimiter

import "testing"

func TestLowPriorityCannotConsumeReservedCapacity(t *testing.T) {
	limiter := New(5, 2)
	for i := 0; i < 3; i++ {
		if !limiter.Acquire(Low) {
			t.Fatalf("low priority acquire %d should succeed", i)
		}
	}
	if limiter.Acquire(Low) {
		t.Fatalf("low priority should stop at non-reserved capacity")
	}
}

func TestHighPriorityCanUseReservedCapacity(t *testing.T) {
	limiter := New(5, 2)
	for i := 0; i < 3; i++ {
		limiter.Acquire(Low)
	}
	for i := 0; i < 2; i++ {
		if !limiter.Acquire(High) {
			t.Fatalf("high priority acquire %d should use reserved capacity", i)
		}
	}
	if limiter.Acquire(High) {
		t.Fatalf("total capacity should be exhausted")
	}
}

func TestReleaseReturnsCapacity(t *testing.T) {
	limiter := New(2, 1)
	if !limiter.Acquire(High) || !limiter.Acquire(Low) {
		t.Fatalf("expected both acquires to succeed")
	}
	limiter.Release(High)
	if !limiter.Acquire(High) {
		t.Fatalf("released high slot should be reusable")
	}
}
