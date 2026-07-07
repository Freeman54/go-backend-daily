package consumerbackpressure

import "testing"

func TestControllerPausesOnHighLag(t *testing.T) {
	controller, err := New(1000, 300, 200, 80)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	paused := controller.Decide(Snapshot{Lag: 1200, InFlight: 30})
	if !paused {
		t.Fatal("Decide() = false, want pause on high lag")
	}
	if !controller.Paused() {
		t.Fatal("controller should remain paused")
	}
}

func TestControllerPausesOnInflightAndUsesHysteresisToResume(t *testing.T) {
	controller, err := New(1000, 300, 200, 80)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if paused := controller.Decide(Snapshot{Lag: 100, InFlight: 220}); !paused {
		t.Fatal("expected pause on in-flight pressure")
	}

	if paused := controller.Decide(Snapshot{Lag: 200, InFlight: 100}); !paused {
		t.Fatal("should stay paused until both metrics recover")
	}

	if paused := controller.Decide(Snapshot{Lag: 280, InFlight: 60}); paused {
		t.Fatal("expected resume when lag and in-flight both recover")
	}
}

func TestControllerRejectsInvalidThresholds(t *testing.T) {
	if _, err := New(100, 200, 10, 5); err == nil {
		t.Fatal("expected invalid lag threshold error")
	}
	if _, err := New(100, 50, 10, 20); err == nil {
		t.Fatal("expected invalid in-flight threshold error")
	}
}
