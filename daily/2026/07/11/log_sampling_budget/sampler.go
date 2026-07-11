package logsamplingbudget

import (
	"sync"
	"time"
)

type clock func() time.Time

type Sampler struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	now    clock
	state  map[string]entry
}

type entry struct {
	windowStart time.Time
	count       int
}

func New(limit int, window time.Duration) *Sampler {
	return &Sampler{
		limit:  limit,
		window: window,
		now:    time.Now,
		state:  make(map[string]entry),
	}
}

func (s *Sampler) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	item := s.state[key]
	if item.windowStart.IsZero() || now.Sub(item.windowStart) >= s.window {
		item.windowStart = now
		item.count = 0
	}

	if item.count >= s.limit {
		s.state[key] = item
		return false
	}

	item.count++
	s.state[key] = item
	return true
}

func (s *Sampler) SetClock(now func() time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.now = now
}
