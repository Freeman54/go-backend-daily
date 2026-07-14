package replicastickiness

import (
	"testing"
	"time"
)

func TestPickKeepsSessionOnSameReplicaWithinTTL(t *testing.T) {
	router := NewRouter(100*time.Millisecond, 2*time.Second)
	base := time.Unix(1000, 0)
	router.now = func() time.Time { return base }

	replicas := []Replica{
		{Name: "replica-a", Lag: 10 * time.Millisecond, Healthy: true},
		{Name: "replica-b", Lag: 20 * time.Millisecond, Healthy: true},
	}

	first := router.Pick("session-1", replicas)
	router.now = func() time.Time { return base.Add(500 * time.Millisecond) }
	second := router.Pick("session-1", []Replica{
		{Name: "replica-b", Lag: 5 * time.Millisecond, Healthy: true},
		{Name: "replica-a", Lag: 10 * time.Millisecond, Healthy: true},
	})

	if first.Replica != "replica-a" {
		t.Fatalf("first replica = %q want replica-a", first.Replica)
	}
	if second.Replica != "replica-a" {
		t.Fatalf("second replica = %q want replica-a", second.Replica)
	}
}

func TestPickMovesSessionWhenPinnedReplicaIsStale(t *testing.T) {
	router := NewRouter(50*time.Millisecond, time.Second)
	base := time.Unix(2000, 0)
	router.now = func() time.Time { return base }

	router.Pick("session-2", []Replica{
		{Name: "replica-a", Lag: 10 * time.Millisecond, Healthy: true},
		{Name: "replica-b", Lag: 20 * time.Millisecond, Healthy: true},
	})

	router.now = func() time.Time { return base.Add(300 * time.Millisecond) }
	route := router.Pick("session-2", []Replica{
		{Name: "replica-a", Lag: 200 * time.Millisecond, Healthy: true},
		{Name: "replica-b", Lag: 10 * time.Millisecond, Healthy: true},
	})

	if route.Replica != "replica-b" {
		t.Fatalf("route replica = %q want replica-b", route.Replica)
	}
}
