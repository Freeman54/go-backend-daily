package batchcheckpoint

import (
	"sync"
	"testing"
)

func TestOutOfOrderAcknowledgements(t *testing.T) {
	tracker := New(10)
	if got := tracker.Ack(12); got != 10 {
		t.Fatalf("checkpoint = %d", got)
	}
	tracker.Ack(10)
	if got := tracker.Ack(11); got != 13 {
		t.Fatalf("checkpoint = %d", got)
	}
	if got := tracker.Ack(9); got != 13 {
		t.Fatalf("old ack changed checkpoint: %d", got)
	}
}
func TestConcurrentAcks(t *testing.T) {
	tracker := New(0)
	var wg sync.WaitGroup
	for i := int64(0); i < 100; i++ {
		wg.Add(1)
		go func(offset int64) { defer wg.Done(); tracker.Ack(offset) }(i)
	}
	wg.Wait()
	if got := tracker.Checkpoint(); got != 100 {
		t.Fatalf("checkpoint = %d", got)
	}
}
