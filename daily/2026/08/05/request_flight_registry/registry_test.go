package request_flight_registry

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRegistryDo_SharesWorkForSameKey(t *testing.T) {
	r := New[string]()
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int32
	work := func(context.Context) (string, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return "value", nil
	}
	result := make(chan string, 2)
	go func() { v, _ := r.Do(context.Background(), "key", work); result <- v }()
	<-started
	go func() { v, _ := r.Do(context.Background(), "key", work); result <- v }()
	time.Sleep(10 * time.Millisecond)
	close(release)
	if <-result != "value" || <-result != "value" {
		t.Fatal("等待者未收到共享结果")
	}
	if calls.Load() != 1 {
		t.Fatalf("work 调用 %d 次，期望 1 次", calls.Load())
	}
}

func TestRegistryDo_CancelsWorkAfterLastWaiterLeaves(t *testing.T) {
	r := New[string]()
	ctx, cancel := context.WithCancel(context.Background())
	workCanceled := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := r.Do(ctx, "key", func(ctx context.Context) (string, error) {
			<-ctx.Done()
			close(workCanceled)
			return "", ctx.Err()
		})
		done <- err
	}()
	cancel()
	if err := <-done; err != context.Canceled {
		t.Fatalf("err = %v", err)
	}
	select {
	case <-workCanceled:
	case <-time.After(time.Second):
		t.Fatal("底层任务未取消")
	}
}
