package perkeymutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLockSerializesSameKey(t *testing.T) {
	locker := New()
	unlock := locker.Lock("order-1")

	entered := make(chan struct{})
	released := make(chan struct{})

	go func() {
		defer close(released)
		unlock2 := locker.Lock("order-1")
		close(entered)
		unlock2()
	}()

	select {
	case <-entered:
		t.Fatal("second goroutine should block on same key")
	case <-time.After(30 * time.Millisecond):
	}

	unlock()

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("second goroutine did not acquire lock after release")
	}

	<-released
}

func TestLockAllowsDifferentKeysInParallel(t *testing.T) {
	locker := New()
	var running int32
	var peak int32
	var wg sync.WaitGroup

	for _, key := range []string{"a", "b"} {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			unlock := locker.Lock(key)
			defer unlock()

			current := atomic.AddInt32(&running, 1)
			for {
				old := atomic.LoadInt32(&peak)
				if current <= old || atomic.CompareAndSwapInt32(&peak, old, current) {
					break
				}
			}

			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&running, -1)
		}(key)
	}

	wg.Wait()

	if peak < 2 {
		t.Fatalf("peak concurrency = %d want at least 2", peak)
	}
}

func TestLockCleansUpIdleKeys(t *testing.T) {
	locker := New()
	unlock := locker.Lock("tenant-42")
	if got := locker.ActiveKeys(); got != 1 {
		t.Fatalf("active keys = %d want 1", got)
	}

	unlock()

	if got := locker.ActiveKeys(); got != 0 {
		t.Fatalf("active keys = %d want 0", got)
	}
}
