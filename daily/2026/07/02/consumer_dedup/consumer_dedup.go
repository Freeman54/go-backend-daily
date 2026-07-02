package consumerdedup

import (
	"context"
	"sync"
	"time"
)

type Message struct {
	ID      string
	Payload string
}

type DedupStore interface {
	MarkOnce(id string, ttl time.Duration) bool
}

type Processor struct {
	Store      DedupStore
	DedupTTL   time.Duration
	HandleFunc func(context.Context, Message) error
}

func (p Processor) Process(ctx context.Context, msg Message) (bool, error) {
	if !p.Store.MarkOnce(msg.ID, p.DedupTTL) {
		return true, nil
	}
	return false, p.HandleFunc(ctx, msg)
}

type InMemoryStore struct {
	mu      sync.Mutex
	now     func() time.Time
	records map[string]time.Time
}

func NewInMemoryStore(now func() time.Time) *InMemoryStore {
	if now == nil {
		now = time.Now
	}
	return &InMemoryStore{
		now:     now,
		records: make(map[string]time.Time),
	}
}

func (s *InMemoryStore) MarkOnce(id string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if expiresAt, ok := s.records[id]; ok && now.Before(expiresAt) {
		return false
	}
	s.records[id] = now.Add(ttl)
	return true
}
