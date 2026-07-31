package orderedshutdowngroup

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestGroup_ShutsDownInDescendingPhaseOrder(t *testing.T) {
	var calls []string
	g := New()
	g.Add(10, "http", func(context.Context) error {
		calls = append(calls, "http")
		return nil
	})
	g.Add(20, "consumer", func(context.Context) error {
		calls = append(calls, "consumer")
		return nil
	})
	g.Add(10, "metrics", func(context.Context) error {
		calls = append(calls, "metrics")
		return nil
	})

	if err := g.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
	if want := []string{"consumer", "http", "metrics"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
}

func TestGroup_CollectsErrorsAndStopsOnCanceledContext(t *testing.T) {
	wantErr := errors.New("flush failed")
	g := New()
	g.Add(20, "flush", func(context.Context) error { return wantErr })
	g.Add(10, "close", func(context.Context) error { return nil })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := g.Shutdown(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Shutdown() error = %v, want context.Canceled", err)
	}
}
