package keyedretryjitter

import (
	"testing"
	"time"
)

func TestDelayGrowsExponentiallyAndCaps(t *testing.T) {
	got := Delay("order-1", 4, time.Second, 5*time.Second, 0)
	if got != 5*time.Second {
		t.Fatalf("expected capped delay, got %v", got)
	}
}

func TestDelayIsStableForSameKeyAndAttempt(t *testing.T) {
	first := Delay("message-42", 3, 2*time.Second, 30*time.Second, 0.4)
	second := Delay("message-42", 3, 2*time.Second, 30*time.Second, 0.4)
	if first != second {
		t.Fatalf("expected stable jitter, got %v and %v", first, second)
	}
}

func TestDelayStaysWithinConfiguredJitterBand(t *testing.T) {
	baseDelay := 8 * time.Second
	got := Delay("message-99", 3, 2*time.Second, 30*time.Second, 0.25)
	min := baseDelay - baseDelay/8
	max := baseDelay + baseDelay/8
	if got < min || got > max {
		t.Fatalf("delay %v outside band [%v, %v]", got, min, max)
	}
}
