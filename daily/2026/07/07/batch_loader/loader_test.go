package batchloader

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestLoaderDeduplicatesKeysWithinBatch(t *testing.T) {
	var (
		mu     sync.Mutex
		calls  [][]int
		loader *Loader[int, string]
		err    error
	)

	loader, err = New[int, string](50*time.Millisecond, 8, func(ctx context.Context, keys []int) (map[int]string, error) {
		mu.Lock()
		calls = append(calls, slices.Clone(keys))
		mu.Unlock()
		return map[int]string{
			1: "user-1",
			2: "user-2",
		}, nil
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	type outcome struct {
		value string
		err   error
	}

	results := make(chan outcome, 3)
	go func() {
		value, err := loader.Load(context.Background(), 1)
		results <- outcome{value: value, err: err}
	}()
	go func() {
		value, err := loader.Load(context.Background(), 1)
		results <- outcome{value: value, err: err}
	}()
	go func() {
		value, err := loader.Load(context.Background(), 2)
		results <- outcome{value: value, err: err}
	}()

	select {
	case <-loader.Wait():
	case <-time.After(time.Second):
		t.Fatal("loader did not signal flush")
	}

	if err := loader.Flush(context.Background()); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	for range 3 {
		outcome := <-results
		if outcome.err != nil {
			t.Fatalf("Load() error = %v", outcome.err)
		}
	}

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 1 {
		t.Fatalf("fetch calls = %d, want 1", len(calls))
	}
	if !sameSet(calls[0], []int{1, 2}) {
		t.Fatalf("fetch keys = %v, want [1 2]", calls[0])
	}
}

func TestLoaderFlushesImmediatelyAtBatchLimit(t *testing.T) {
	loader, err := New[int, string](time.Hour, 2, func(ctx context.Context, keys []int) (map[int]string, error) {
		return map[int]string{
			10: "a",
			20: "b",
		}, nil
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		value, err := loader.Load(context.Background(), 10)
		if err != nil || value != "a" {
			t.Errorf("first Load() = (%q, %v)", value, err)
		}
	}()

	go func() {
		value, err := loader.Load(context.Background(), 20)
		if err != nil || value != "b" {
			t.Errorf("second Load() = (%q, %v)", value, err)
		}
	}()

	select {
	case <-loader.Wait():
	case <-time.After(time.Second):
		t.Fatal("loader did not flush at batch limit")
	}

	if err := loader.Flush(context.Background()); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}
	<-done
}

func TestLoaderPropagatesFetchError(t *testing.T) {
	wantErr := errors.New("db unavailable")
	loader, err := New[int, string](20*time.Millisecond, 4, func(context.Context, []int) (map[int]string, error) {
		return nil, wantErr
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	resultCh := make(chan error, 1)
	go func() {
		_, err := loader.Load(context.Background(), 42)
		resultCh <- err
	}()

	select {
	case <-loader.Wait():
	case <-time.After(time.Second):
		t.Fatal("loader did not signal flush")
	}

	err = loader.Flush(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("Flush() error = %v, want %v", err, wantErr)
	}
	if got := <-resultCh; !errors.Is(got, wantErr) {
		t.Fatalf("Load() error = %v, want %v", got, wantErr)
	}
}

func sameSet(got, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	got = slices.Clone(got)
	want = slices.Clone(want)
	slices.Sort(got)
	slices.Sort(want)
	return slices.Equal(got, want)
}
