package cache_write_coalescer

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCoalescer_DoSharesConcurrentWrite(t *testing.T) {
	c := New()
	start := make(chan struct{})
	var calls atomic.Int32
	fn := func() error {
		calls.Add(1)
		<-start
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(2)
	for range 2 {
		go func() {
			defer wg.Done()
			if err := c.Do("user:7", fn); err != nil {
				t.Errorf("Do() error = %v", err)
			}
		}()
	}
	time.Sleep(10 * time.Millisecond)
	close(start)
	wg.Wait()
	if got := calls.Load(); got != 1 {
		t.Fatalf("write calls = %d, want 1", got)
	}
}

func TestCoalescer_DoDoesNotMergeDifferentKeys(t *testing.T) {
	c := New()
	var calls atomic.Int32
	if err := c.Do("a", func() error { calls.Add(1); return nil }); err != nil {
		t.Fatal(err)
	}
	if err := c.Do("b", func() error { calls.Add(1); return nil }); err != nil {
		t.Fatal(err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("write calls = %d, want 2", got)
	}
}
