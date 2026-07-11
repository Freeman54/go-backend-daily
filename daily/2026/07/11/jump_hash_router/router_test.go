package jumphashrouter

import (
	"fmt"
	"testing"
)

func TestRouteIsStableForSameInput(t *testing.T) {
	router := New(16)
	got := router.Route("order-1024")
	for i := 0; i < 10; i++ {
		if again := router.Route("order-1024"); again != got {
			t.Fatalf("route changed from %d to %d", got, again)
		}
	}
}

func TestMovedRatioStaysLowWhenAddingOnePartition(t *testing.T) {
	keys := make([]string, 0, 2000)
	for i := 0; i < 2000; i++ {
		keys = append(keys, fmt.Sprintf("user-%d", i))
	}

	ratio := MovedRatio(keys, 8, 9)
	if ratio <= 0 || ratio >= 0.25 {
		t.Fatalf("unexpected moved ratio: %f", ratio)
	}
}

func TestNewWithInvalidPartitionCountReturnsNegativeRoute(t *testing.T) {
	if got := New(0).Route("any"); got != -1 {
		t.Fatalf("route = %d, want -1", got)
	}
}
