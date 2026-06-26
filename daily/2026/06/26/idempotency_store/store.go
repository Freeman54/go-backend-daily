package idempotencystore

import (
	"sync"
	"time"
)

type Response struct {
	StatusCode int
	Body       string
}

type entry struct {
	response  Response
	expiresAt time.Time
}

type Store struct {
	mu      sync.Mutex
	entries map[string]entry
	now     func() time.Time
}

func NewStore() *Store {
	return &Store{
		entries: make(map[string]entry),
		now:     time.Now,
	}
}

func (s *Store) Remember(key string, response Response, ttl time.Duration) {
	s.RememberAt(key, response, ttl, s.now())
}

func (s *Store) RememberAt(key string, response Response, ttl time.Duration, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[key] = entry{
		response:  response,
		expiresAt: now.Add(ttl),
	}
}

func (s *Store) Lookup(key string, now time.Time) (Response, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.entries[key]
	if !ok {
		return Response{}, false
	}
	if !now.Before(item.expiresAt) {
		delete(s.entries, key)
		return Response{}, false
	}
	return item.response, true
}
