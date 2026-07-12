package batchflusher

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestFlusherFlushesOnBatchSize(t *testing.T) {
	var flushed [][]int
	flusher := New(3, time.Minute, func(items []int) error {
		flushed = append(flushed, items)
		return nil
	})

	now := time.Unix(0, 0)
	for _, item := range []int{1, 2, 3} {
		if err := flusher.Add(item, now); err != nil {
			t.Fatalf("add failed: %v", err)
		}
	}

	want := [][]int{{1, 2, 3}}
	if !reflect.DeepEqual(flushed, want) {
		t.Fatalf("flushed = %#v want %#v", flushed, want)
	}
}

func TestFlusherFlushesOnMaxWait(t *testing.T) {
	var flushed [][]string
	flusher := New(10, 5*time.Second, func(items []string) error {
		flushed = append(flushed, items)
		return nil
	})

	base := time.Unix(0, 0)
	if err := flusher.Add("a", base); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if err := flusher.Add("b", base.Add(time.Second)); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	if err := flusher.FlushDue(base.Add(6 * time.Second)); err != nil {
		t.Fatalf("flush due failed: %v", err)
	}

	want := [][]string{{"a", "b"}}
	if !reflect.DeepEqual(flushed, want) {
		t.Fatalf("flushed = %#v want %#v", flushed, want)
	}
}

func TestFlusherKeepsBatchWhenFlushFails(t *testing.T) {
	wantErr := errors.New("db down")
	calls := 0
	flusher := New(2, time.Second, func(items []int) error {
		calls++
		if calls == 1 {
			return wantErr
		}
		return nil
	})

	now := time.Unix(0, 0)
	if err := flusher.Add(1, now); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	err := flusher.Add(2, now)
	if !errors.Is(err, wantErr) {
		t.Fatalf("flush error = %v want %v", err, wantErr)
	}
	if err := flusher.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}
