package retryafterhint

import (
	"testing"
	"time"
)

func TestSecondsRoundsUpSubSecondDelay(t *testing.T) {
	now := time.Unix(100, 0)
	got := Seconds(now, now.Add(1500*time.Millisecond), 30)
	if got != "2" {
		t.Fatalf("expected 2, got %s", got)
	}
}

func TestSecondsKeepsMinimumOneSecond(t *testing.T) {
	now := time.Unix(100, 0)
	got := Seconds(now, now.Add(-time.Second), 30)
	if got != "1" {
		t.Fatalf("expected 1, got %s", got)
	}
}

func TestSecondsClampsToMax(t *testing.T) {
	now := time.Unix(100, 0)
	got := Seconds(now, now.Add(3*time.Minute), 15)
	if got != "15" {
		t.Fatalf("expected 15, got %s", got)
	}
}
