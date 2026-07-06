package replicalagguard

import (
	"testing"
	"time"
)

func TestChoosePrimaryForReadAfterWrite(t *testing.T) {
	decision := Choose(true, 50*time.Millisecond, []Replica{
		{Name: "replica-a", Lag: 5 * time.Millisecond, Healthy: true},
	})

	if decision.Target != "primary" || decision.Reason != "read-after-write" {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestChooseFastestHealthyReplicaWithinLag(t *testing.T) {
	decision := Choose(false, 50*time.Millisecond, []Replica{
		{Name: "replica-a", Lag: 40 * time.Millisecond, Healthy: true},
		{Name: "replica-b", Lag: 10 * time.Millisecond, Healthy: true},
		{Name: "replica-c", Lag: 5 * time.Millisecond, Healthy: false},
	})

	if decision.Target != "replica-b" || decision.Reason != "replica-ok" {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestFallbackToPrimaryWhenLagTooHigh(t *testing.T) {
	decision := Choose(false, 20*time.Millisecond, []Replica{
		{Name: "replica-a", Lag: 25 * time.Millisecond, Healthy: true},
	})

	if decision.Target != "primary" || decision.Reason != "replica-lag" {
		t.Fatalf("decision = %+v", decision)
	}
}
