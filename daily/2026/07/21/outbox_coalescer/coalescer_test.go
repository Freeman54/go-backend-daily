package outboxcoalescer

import "testing"

func TestCoalescerKeepsLatestVersionPerAggregate(t *testing.T) {
	c := New()
	mustAdd(t, c, Event{AggregateID: "order-2", Version: 1, Kind: "created"})
	mustAdd(t, c, Event{AggregateID: "order-1", Version: 2, Kind: "paid"})
	mustAdd(t, c, Event{AggregateID: "order-1", Version: 1, Kind: "created"})
	mustAdd(t, c, Event{AggregateID: "order-2", Version: 3, Kind: "cancelled"})

	got := c.Flush()
	want := []Event{
		{AggregateID: "order-1", Version: 2, Kind: "paid"},
		{AggregateID: "order-2", Version: 3, Kind: "cancelled"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d events, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("event %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestCoalescerRejectsInvalidEvent(t *testing.T) {
	c := New()
	if err := c.Add(Event{}); err == nil {
		t.Fatal("expected invalid event to fail")
	}
}

func mustAdd(t *testing.T, c *Coalescer, event Event) {
	t.Helper()
	if err := c.Add(event); err != nil {
		t.Fatalf("add %#v: %v", event, err)
	}
}
