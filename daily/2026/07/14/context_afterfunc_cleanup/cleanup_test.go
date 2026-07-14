package contextafterfunccleanup

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestCleanupRunsInReverseOrderOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var group Group
	var mu sync.Mutex
	var order []string

	group.Add(func() {
		mu.Lock()
		order = append(order, "conn")
		mu.Unlock()
	})
	group.Add(func() {
		mu.Lock()
		order = append(order, "span")
		mu.Unlock()
	})

	group.Bind(ctx)
	cancel()

	group.Cleanup()

	mu.Lock()
	defer mu.Unlock()
	want := []string{"span", "conn"}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("cleanup order = %v want %v", order, want)
	}
}

func TestStopPreventsCleanupUntilManualCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var group Group
	called := 0
	group.Add(func() { called++ })

	stop := group.Bind(ctx)
	if !stop() {
		t.Fatalf("stop should detach cleanup callback")
	}

	cancel()
	if called != 0 {
		t.Fatalf("cleanup called = %d want 0", called)
	}

	group.Cleanup()
	if called != 1 {
		t.Fatalf("cleanup called = %d want 1", called)
	}
}
