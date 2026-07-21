package outboxcoalescer

import (
	"fmt"
	"slices"
)

type Event struct {
	AggregateID string
	Version     int64
	Kind        string
}

type Coalescer struct {
	latest map[string]Event
}

func New() *Coalescer {
	return &Coalescer{
		latest: make(map[string]Event),
	}
}

func (c *Coalescer) Add(event Event) error {
	if event.AggregateID == "" {
		return fmt.Errorf("aggregate id must not be empty")
	}
	if event.Version <= 0 {
		return fmt.Errorf("version must be positive")
	}
	if event.Kind == "" {
		return fmt.Errorf("kind must not be empty")
	}

	current, ok := c.latest[event.AggregateID]
	if !ok || event.Version >= current.Version {
		c.latest[event.AggregateID] = event
	}
	return nil
}

func (c *Coalescer) Flush() []Event {
	ids := make([]string, 0, len(c.latest))
	for id := range c.latest {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	events := make([]Event, 0, len(ids))
	for _, id := range ids {
		events = append(events, c.latest[id])
	}
	return events
}
