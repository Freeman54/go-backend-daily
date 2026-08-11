package generationlatch

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLatchWakesAllWaiters(t *testing.T) {
	latch := New()
	results := make(chan uint64, 2)
	for range 2 {
		go func() {
			generation, err := latch.Wait(context.Background(), 0)
			if err == nil {
				results <- generation
			}
		}()
	}
	if got := latch.Advance(); got != 1 {
		t.Fatalf("Advance() = %d", got)
	}
	for range 2 {
		select {
		case got := <-results:
			if got != 1 {
				t.Fatalf("generation = %d", got)
			}
		case <-time.After(time.Second):
			t.Fatal("waiter was not awakened")
		}
	}
}

func TestLatchAvoidsLostWakeupAndHonorsContext(t *testing.T) {
	latch := New()
	latch.Advance()
	if got, err := latch.Wait(context.Background(), 0); err != nil || got != 1 {
		t.Fatalf("Wait() = %d, %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := latch.Wait(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
}
