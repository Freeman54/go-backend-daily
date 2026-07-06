package lagautoscaler

import "testing"

func TestDecideScalesUpWithinStepLimit(t *testing.T) {
	policy := Policy{
		TargetCatchupSeconds: 60,
		MinConsumers:         1,
		MaxConsumers:         10,
		MaxScaleStep:         2,
	}
	snapshot := Snapshot{
		Lag:                   2400,
		InflowPerSecond:       10,
		CurrentConsumers:      2,
		ThroughputPerConsumer: 10,
	}

	got := Decide(policy, snapshot)
	if got != 4 {
		t.Fatalf("Decide() = %d, want 4", got)
	}
}

func TestDecideUsesInflowWhenLagIsLow(t *testing.T) {
	policy := Policy{
		TargetCatchupSeconds: 60,
		MinConsumers:         1,
		MaxConsumers:         10,
		MaxScaleStep:         3,
	}
	snapshot := Snapshot{
		Lag:                   0,
		InflowPerSecond:       45,
		CurrentConsumers:      6,
		ThroughputPerConsumer: 10,
	}

	got := Decide(policy, snapshot)
	if got != 5 {
		t.Fatalf("Decide() = %d, want 5", got)
	}
}

func TestDecideClampsToMinimumWhenMetricsAreZero(t *testing.T) {
	policy := Policy{
		TargetCatchupSeconds: 60,
		MinConsumers:         2,
		MaxConsumers:         10,
		MaxScaleStep:         1,
	}
	snapshot := Snapshot{
		Lag:                   0,
		InflowPerSecond:       0,
		CurrentConsumers:      4,
		ThroughputPerConsumer: 0,
	}

	got := Decide(policy, snapshot)
	if got != 4 {
		t.Fatalf("Decide() = %d, want 4", got)
	}
}
