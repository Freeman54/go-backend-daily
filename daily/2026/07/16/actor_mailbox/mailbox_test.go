package actormailbox

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSubmitPreservesPerKeyOrder(t *testing.T) {
	mailbox := New(8)
	var (
		mu    sync.Mutex
		order []int
		wg    sync.WaitGroup
	)

	for i := 1; i <= 3; i++ {
		i := i
		wg.Add(1)
		err := mailbox.Submit(context.Background(), "order-1", func(context.Context) {
			defer wg.Done()
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
		})
		if err != nil {
			t.Fatalf("submit: %v", err)
		}
	}

	wg.Wait()

	if !reflect.DeepEqual(order, []int{1, 2, 3}) {
		t.Fatalf("unexpected order: %#v", order)
	}
}

func TestSubmitRunsDifferentKeysConcurrently(t *testing.T) {
	mailbox := New(4)
	var wg sync.WaitGroup
	started := make(chan string, 2)
	release := make(chan struct{})

	submit := func(key string) {
		wg.Add(1)
		err := mailbox.Submit(context.Background(), key, func(context.Context) {
			defer wg.Done()
			started <- key
			<-release
		})
		if err != nil {
			t.Fatalf("submit %s: %v", key, err)
		}
	}

	submit("a")
	submit("b")

	got := map[string]bool{
		<-started: true,
		<-started: true,
	}
	close(release)
	wg.Wait()

	if len(got) != 2 || !got["a"] || !got["b"] {
		t.Fatalf("tasks did not run concurrently across keys: %#v", got)
	}
}

func TestSubmitReturnsQueueFull(t *testing.T) {
	mailbox := New(1)
	block := make(chan struct{})
	err := mailbox.Submit(context.Background(), "tenant", func(context.Context) {
		<-block
	})
	if err != nil {
		t.Fatalf("submit first task: %v", err)
	}

	err = mailbox.Submit(context.Background(), "tenant", func(context.Context) {})
	if err != nil {
		t.Fatalf("submit queued task: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		err = mailbox.Submit(context.Background(), "tenant", func(context.Context) {})
		if errors.Is(err, ErrQueueFull) {
			close(block)
			return
		}
		if time.Now().After(deadline) {
			close(block)
			t.Fatal("expected ErrQueueFull")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
