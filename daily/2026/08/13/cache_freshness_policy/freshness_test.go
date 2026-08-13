package freshnesspolicy

import (
	"testing"
	"time"
)

func TestClassify(t *testing.T) {
	base := time.Unix(100, 0)
	p := Policy{FreshFor: 10 * time.Second, ServeStaleFor: 5 * time.Second}
	for _, tt := range []struct {
		age  time.Duration
		want State
	}{{5 * time.Second, Fresh}, {10 * time.Second, Fresh}, {12 * time.Second, Stale}, {15 * time.Second, Stale}, {16 * time.Second, Expired}, {-time.Second, Fresh}} {
		if got := p.Classify(base, base.Add(tt.age)); got != tt.want {
			t.Fatalf("age %v: got %v want %v", tt.age, got, tt.want)
		}
	}
}

func TestNegativeDurations(t *testing.T) {
	if got := (Policy{FreshFor: -1, ServeStaleFor: -1}).Classify(time.Unix(0, 0), time.Unix(1, 0)); got != Expired {
		t.Fatalf("got %v", got)
	}
}
