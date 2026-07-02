package consumerdedup

import (
	"context"
	"testing"
	"time"
)

func TestProcessSkipsDuplicateMessage(t *testing.T) {
	t.Parallel()

	store := NewInMemoryStore(func() time.Time { return time.Unix(0, 0) })
	calls := 0
	processor := Processor{
		Store:    store,
		DedupTTL: time.Minute,
		HandleFunc: func(context.Context, Message) error {
			calls++
			return nil
		},
	}

	duplicate, err := processor.Process(context.Background(), Message{ID: "msg-1"})
	if err != nil || duplicate {
		t.Fatalf("first Process() duplicate=%v err=%v, want false nil", duplicate, err)
	}
	duplicate, err = processor.Process(context.Background(), Message{ID: "msg-1"})
	if err != nil || !duplicate {
		t.Fatalf("second Process() duplicate=%v err=%v, want true nil", duplicate, err)
	}
	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
}

func TestProcessAllowsRetryAfterTTL(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	store := NewInMemoryStore(func() time.Time { return now })
	calls := 0
	processor := Processor{
		Store:    store,
		DedupTTL: time.Minute,
		HandleFunc: func(context.Context, Message) error {
			calls++
			return nil
		},
	}

	if _, err := processor.Process(context.Background(), Message{ID: "msg-1"}); err != nil {
		t.Fatalf("first Process() error = %v", err)
	}
	now = now.Add(2 * time.Minute)
	duplicate, err := processor.Process(context.Background(), Message{ID: "msg-1"})
	if err != nil || duplicate {
		t.Fatalf("second Process() duplicate=%v err=%v, want false nil", duplicate, err)
	}
	if calls != 2 {
		t.Fatalf("handler calls = %d, want 2", calls)
	}
}
