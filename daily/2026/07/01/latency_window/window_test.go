package latencywindow

import (
	"testing"
	"time"
)

func TestWindowEvictsExpiredSamples(t *testing.T) {
	t.Parallel()

	base := time.Unix(0, 0)
	window := NewWindow(time.Minute, 100*time.Millisecond)
	window.Record(base, 90*time.Millisecond)
	window.Record(base.Add(30*time.Second), 120*time.Millisecond)
	window.Record(base.Add(90*time.Second), 80*time.Millisecond)

	if got := window.Count(base.Add(90 * time.Second)); got != 2 {
		t.Fatalf("Count() = %d, want 2", got)
	}
}

func TestWindowComputesSuccessRatio(t *testing.T) {
	t.Parallel()

	base := time.Unix(0, 0)
	window := NewWindow(time.Minute, 100*time.Millisecond)
	window.Record(base, 80*time.Millisecond)
	window.Record(base.Add(time.Second), 90*time.Millisecond)
	window.Record(base.Add(2*time.Second), 150*time.Millisecond)

	if got := window.SuccessRatio(base.Add(2 * time.Second)); got != 2.0/3.0 {
		t.Fatalf("SuccessRatio() = %v, want %v", got, 2.0/3.0)
	}
}

func TestWindowComputesP95(t *testing.T) {
	t.Parallel()

	base := time.Unix(0, 0)
	window := NewWindow(time.Minute, 200*time.Millisecond)
	for _, latency := range []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
		60 * time.Millisecond,
		70 * time.Millisecond,
		80 * time.Millisecond,
		90 * time.Millisecond,
		300 * time.Millisecond,
	} {
		window.Record(base, latency)
	}

	if got := window.P95(base); got != 300*time.Millisecond {
		t.Fatalf("P95() = %v, want 300ms", got)
	}
}
