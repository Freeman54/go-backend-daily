package adaptiveconcurrency

import (
	"testing"
	"time"
)

func TestLimiterRampsUpAfterStableLowLatency(t *testing.T) {
	limiter := New(2, 6, 100*time.Millisecond, 2)

	for i := 0; i < 2; i++ {
		limiter.Observe(Sample{Latency: 30 * time.Millisecond})
	}
	if got := limiter.Limit(); got != 3 {
		t.Fatalf("expected limit to grow to 3, got %d", got)
	}

	for i := 0; i < 4; i++ {
		limiter.Observe(Sample{Latency: 30 * time.Millisecond})
	}
	if got := limiter.Limit(); got != 5 {
		t.Fatalf("expected limit to continue ramping to 5, got %d", got)
	}
}

func TestLimiterCutsLimitOnHighLatency(t *testing.T) {
	limiter := New(2, 8, 100*time.Millisecond, 1)

	for i := 0; i < 4; i++ {
		limiter.Observe(Sample{Latency: 20 * time.Millisecond})
	}
	if got := limiter.Limit(); got != 6 {
		t.Fatalf("expected limit to ramp to 6, got %d", got)
	}

	limiter.Observe(Sample{Latency: 150 * time.Millisecond})
	if got := limiter.Limit(); got != 3 {
		t.Fatalf("expected limit to halve to 3, got %d", got)
	}
}

func TestLimiterCutsLimitOnFailures(t *testing.T) {
	limiter := New(1, 4, 100*time.Millisecond, 1)

	for i := 0; i < 3; i++ {
		limiter.Observe(Sample{Latency: 10 * time.Millisecond})
	}
	if got := limiter.Limit(); got != 4 {
		t.Fatalf("expected limit to reach max 4, got %d", got)
	}

	limiter.Observe(Sample{Latency: 10 * time.Millisecond, Failed: true})
	if got := limiter.Limit(); got != 2 {
		t.Fatalf("expected failure to reduce limit to 2, got %d", got)
	}
}
