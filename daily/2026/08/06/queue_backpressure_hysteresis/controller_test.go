package queue_backpressure_hysteresis

import "testing"

func TestController_HoldsSheddingUntilLowWatermark(t *testing.T) {
	c, err := New(3, 7)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		depth int
		want  bool
	}{{6, false}, {7, true}, {5, true}, {3, false}} {
		if got := c.Observe(tc.depth); got != tc.want {
			t.Fatalf("Observe(%d) = %v, want %v", tc.depth, got, tc.want)
		}
	}
}

func TestNew_RejectsInvalidWatermarks(t *testing.T) {
	if _, err := New(8, 7); err == nil {
		t.Fatal("New() should require low < high")
	}
}
