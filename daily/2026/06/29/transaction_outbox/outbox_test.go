package transactionoutbox

import (
	"context"
	"errors"
	"testing"
)

type stubPublisher struct {
	failOnce  bool
	published []Event
}

func (s *stubPublisher) Publish(_ context.Context, event Event) error {
	if s.failOnce {
		s.failOnce = false
		return errors.New("broker unavailable")
	}
	s.published = append(s.published, event)
	return nil
}

func TestSaveOrderCreatesPendingEvent(t *testing.T) {
	t.Parallel()

	store := NewStore()
	if err := store.SaveOrder(Order{ID: "order-1", Amount: 99}); err != nil {
		t.Fatalf("SaveOrder() error = %v", err)
	}

	pending := store.PendingEvents()
	if len(pending) != 1 {
		t.Fatalf("pending events = %d, want 1", len(pending))
	}
	if pending[0].Payload != "order-1" {
		t.Fatalf("event payload = %q, want %q", pending[0].Payload, "order-1")
	}
}

func TestDispatcherRetriesUnpublishedEvent(t *testing.T) {
	t.Parallel()

	store := NewStore()
	if err := store.SaveOrder(Order{ID: "order-2", Amount: 199}); err != nil {
		t.Fatalf("SaveOrder() error = %v", err)
	}

	publisher := &stubPublisher{failOnce: true}
	dispatcher := NewDispatcher(store, publisher)

	if err := dispatcher.FlushOnce(context.Background()); err == nil {
		t.Fatal("FlushOnce() error = nil, want publish failure")
	}
	if len(store.PendingEvents()) != 1 {
		t.Fatalf("pending events after failure = %d, want 1", len(store.PendingEvents()))
	}

	if err := dispatcher.FlushOnce(context.Background()); err != nil {
		t.Fatalf("FlushOnce() retry error = %v", err)
	}
	if len(store.PendingEvents()) != 0 {
		t.Fatalf("pending events after retry = %d, want 0", len(store.PendingEvents()))
	}
	if len(publisher.published) != 1 {
		t.Fatalf("published count = %d, want 1", len(publisher.published))
	}
}
