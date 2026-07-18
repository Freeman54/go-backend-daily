package writebehindcache

import "sync"

type Update struct {
	Key   string
	Value string
}

type Store struct {
	mu    sync.Mutex
	items map[string]string
	queue []Update
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]string),
		queue: make([]Update, 0),
	}
}

func (s *Store) Enqueue(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, Update{Key: key, Value: value})
}

func (s *Store) Flush(limit int) []Update {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit <= 0 || limit > len(s.queue) {
		limit = len(s.queue)
	}
	batch := append([]Update(nil), s.queue[:limit]...)
	s.queue = append([]Update(nil), s.queue[limit:]...)

	for _, item := range batch {
		s.items[item.Key] = item.Value
	}
	return batch
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.items[key]
	return value, ok
}

func (s *Store) Pending() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.queue)
}
