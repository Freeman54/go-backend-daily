package visibilityheartbeat

import (
	"testing"
	"time"
)

func TestLeaseOnlyRenewsNearExpiry(t *testing.T) {
	start := time.Unix(0, 0)
	lease := NewLease(start, 30*time.Second, 10*time.Second, 90*time.Second)

	if lease.Due(start.Add(15 * time.Second)) {
		t.Fatal("lease should not be due too early")
	}
	if !lease.Due(start.Add(20 * time.Second)) {
		t.Fatal("lease should be due inside renew window")
	}
}

func TestLeaseHeartbeatExtendsVisibility(t *testing.T) {
	start := time.Unix(0, 0)
	lease := NewLease(start, 30*time.Second, 10*time.Second, 90*time.Second)

	renewed, giveUp := lease.Heartbeat(start.Add(25 * time.Second))
	if !renewed || giveUp {
		t.Fatalf("expected renewal, renewed=%v giveUp=%v", renewed, giveUp)
	}
	if got := lease.ExpiresAt(); !got.Equal(start.Add(60 * time.Second)) {
		t.Fatalf("expected expiry at 60s, got %v", got.Sub(start))
	}
}

func TestLeaseStopsRenewingAfterCap(t *testing.T) {
	start := time.Unix(0, 0)
	lease := NewLease(start, 30*time.Second, 10*time.Second, 60*time.Second)

	if renewed, giveUp := lease.Heartbeat(start.Add(25 * time.Second)); !renewed || giveUp {
		t.Fatalf("expected first renewal, renewed=%v giveUp=%v", renewed, giveUp)
	}
	if renewed, giveUp := lease.Heartbeat(start.Add(55 * time.Second)); renewed || !giveUp {
		t.Fatalf("expected give up at extension cap, renewed=%v giveUp=%v", renewed, giveUp)
	}
}
