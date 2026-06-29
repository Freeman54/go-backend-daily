package transactionoutbox

import (
	"context"
	"fmt"
	"sync"
)

type Order struct {
	ID     string
	Amount int64
}

type Event struct {
	ID        int64
	Topic     string
	Payload   string
	Published bool
}

type Store struct {
	mu          sync.Mutex
	orders      map[string]Order
	events      []Event
	nextEventID int64
}

func NewStore() *Store {
	return &Store{
		orders: make(map[string]Order),
	}
}

func (s *Store) SaveOrder(order Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[order.ID]; exists {
		return fmt.Errorf("order %s already exists", order.ID)
	}

	s.orders[order.ID] = order
	s.nextEventID++
	s.events = append(s.events, Event{
		ID:      s.nextEventID,
		Topic:   "order.created",
		Payload: order.ID,
	})
	return nil
}

func (s *Store) PendingEvents() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	var pending []Event
	for _, event := range s.events {
		if !event.Published {
			pending = append(pending, event)
		}
	}
	return pending
}

func (s *Store) MarkPublished(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.events {
		if s.events[i].ID == id {
			s.events[i].Published = true
			return
		}
	}
}

type Publisher interface {
	Publish(context.Context, Event) error
}

type Dispatcher struct {
	store     *Store
	publisher Publisher
}

func NewDispatcher(store *Store, publisher Publisher) *Dispatcher {
	return &Dispatcher{store: store, publisher: publisher}
}

func (d *Dispatcher) FlushOnce(ctx context.Context) error {
	for _, event := range d.store.PendingEvents() {
		if err := d.publisher.Publish(ctx, event); err != nil {
			return err
		}
		d.store.MarkPublished(event.ID)
	}
	return nil
}
