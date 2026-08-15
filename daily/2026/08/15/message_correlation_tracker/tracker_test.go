package messagecorrelationtracker

import (
	"errors"
	"sync"
	"testing"
)

func TestRegisterAndResolve(t *testing.T) {
	tracker := New()
	ch, err := tracker.Register("req-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := tracker.Resolve("req-1", Result{Payload: []byte("ok")}); err != nil {
		t.Fatal(err)
	}
	got := <-ch
	if string(got.Payload) != "ok" || tracker.Pending() != 0 {
		t.Fatalf("result = %#v, pending = %d", got, tracker.Pending())
	}
}

func TestDuplicateUnknownAndCancel(t *testing.T) {
	tracker := New()
	if _, err := tracker.Register(""); !errors.Is(err, ErrUnknown) {
		t.Fatalf("empty id error = %v", err)
	}
	if _, err := tracker.Register("req-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := tracker.Register("req-1"); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate error = %v", err)
	}
	if !tracker.Cancel("req-1") || tracker.Cancel("req-1") {
		t.Fatal("unexpected cancel result")
	}
	if err := tracker.Resolve("missing", Result{}); !errors.Is(err, ErrUnknown) {
		t.Fatalf("unknown error = %v", err)
	}
}

func TestConcurrentResolveOnlySucceedsOnce(t *testing.T) {
	tracker := New()
	ch, _ := tracker.Register("req-1")
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); errs <- tracker.Resolve("req-1", Result{}) }()
	}
	wg.Wait()
	close(errs)
	success := 0
	for err := range errs {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("successes = %d", success)
	}
	<-ch
}
