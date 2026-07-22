package hedgedrequestgate

import (
	"testing"
	"time"
)

func TestGateAllowsHedgeAfterDelayWithinLimit(t *testing.T) {
	g, err := New(50*time.Millisecond, 2)
	if err != nil {
		t.Fatalf("new gate: %v", err)
	}
	now := time.Unix(0, 0)
	request := g.Start(now)
	if request.ShouldHedge(now.Add(49 * time.Millisecond)) {
		t.Fatal("hedged too early")
	}
	if !request.ShouldHedge(now.Add(50 * time.Millisecond)) {
		t.Fatal("expected hedge")
	}
	if !request.ShouldHedge(now.Add(100 * time.Millisecond)) {
		t.Fatal("expected second hedge")
	}
	if request.ShouldHedge(now.Add(150 * time.Millisecond)) {
		t.Fatal("exceeded hedge limit")
	}
}

func TestNewRejectsInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct {
		delay time.Duration
		max   int
	}{{0, 1}, {time.Second, 0}} {
		if _, err := New(tc.delay, tc.max); err == nil {
			t.Fatalf("expected error for %#v", tc)
		}
	}
}
