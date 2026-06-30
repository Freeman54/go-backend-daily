package readinessgate

import (
	"reflect"
	"testing"
)

func TestSnapshotReadyWhenAllDependenciesPass(t *testing.T) {
	t.Parallel()

	gate := New("db", "kafka")
	gate.SetDependency("db", true)
	gate.SetDependency("kafka", true)

	got := gate.Snapshot()
	if !got.Ready {
		t.Fatalf("Snapshot().Ready = false, want true")
	}
	if len(got.Reasons) != 0 {
		t.Fatalf("Snapshot().Reasons = %v, want empty", got.Reasons)
	}
}

func TestSnapshotListsBlockingReasons(t *testing.T) {
	t.Parallel()

	gate := New("db", "redis")
	gate.SetDependency("db", true)
	gate.SetDraining(true)

	got := gate.Snapshot()
	want := []string{"instance is draining", "redis is not ready"}
	if got.Ready {
		t.Fatalf("Snapshot().Ready = true, want false")
	}
	if !reflect.DeepEqual(got.Reasons, want) {
		t.Fatalf("Snapshot().Reasons = %v, want %v", got.Reasons, want)
	}
}
