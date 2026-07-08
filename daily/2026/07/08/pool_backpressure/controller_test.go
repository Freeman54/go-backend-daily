package poolbackpressure

import "testing"

func TestDecideAdmitWhenPoolIsHealthy(t *testing.T) {
	controller := New(0.7, 0.9, 10)

	decision := controller.Decide(Snapshot{
		InUse:   5,
		MaxOpen: 20,
		Waiters: 1,
	})
	if decision != DecisionAdmit {
		t.Fatalf("expected admit, got %s", decision)
	}
}

func TestDecideDegradeWhenUsageIsHigh(t *testing.T) {
	controller := New(0.7, 0.9, 10)

	decision := controller.Decide(Snapshot{
		InUse:   15,
		MaxOpen: 20,
		Waiters: 2,
	})
	if decision != DecisionDegrade {
		t.Fatalf("expected degrade, got %s", decision)
	}
}

func TestDecideRejectWhenWaitersExplode(t *testing.T) {
	controller := New(0.7, 0.9, 10)

	decision := controller.Decide(Snapshot{
		InUse:   12,
		MaxOpen: 20,
		Waiters: 12,
	})
	if decision != DecisionReject {
		t.Fatalf("expected reject, got %s", decision)
	}
}
