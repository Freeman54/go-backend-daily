package writebehindbuffer

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBufferFlushesMergedMutations(t *testing.T) {
	buffer := New(2)
	if flushNow := buffer.Add(Mutation{Key: "a", Delta: 1}); flushNow {
		t.Fatalf("first key should not cross threshold")
	}
	if flushNow := buffer.Add(Mutation{Key: "a", Delta: 3}); flushNow {
		t.Fatalf("same key should still count as one pending key")
	}
	if flushNow := buffer.Add(Mutation{Key: "b", Delta: 2}); !flushNow {
		t.Fatalf("second distinct key should cross threshold")
	}

	got := buffer.Flush()
	if got["a"] != 4 || got["b"] != 2 {
		t.Fatalf("unexpected flush payload: %#v", got)
	}
}

func TestBufferRunFlushesOnIntervalAndShutdown(t *testing.T) {
	buffer := New(10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu      sync.Mutex
		batches []map[string]int
	)

	done := make(chan struct{})
	go func() {
		buffer.Run(ctx, 10*time.Millisecond, func(batch map[string]int) {
			mu.Lock()
			batches = append(batches, batch)
			mu.Unlock()
		})
		close(done)
	}()

	buffer.Add(Mutation{Key: "a", Delta: 1})
	time.Sleep(25 * time.Millisecond)
	buffer.Add(Mutation{Key: "b", Delta: 2})
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if len(batches) < 2 {
		t.Fatalf("expected interval flush and shutdown flush, got %d batches", len(batches))
	}
	if batches[0]["a"] != 1 {
		t.Fatalf("expected first batch to flush key a, got %#v", batches[0])
	}
	last := batches[len(batches)-1]
	if last["b"] != 2 {
		t.Fatalf("expected final batch to flush key b, got %#v", last)
	}
}
