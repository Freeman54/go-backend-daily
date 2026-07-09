package versionedcachekeys

import (
	"fmt"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	versions map[string]int
	values   map[string]string
}

func New() *Store {
	return &Store{
		versions: make(map[string]int),
		values:   make(map[string]string),
	}
}

func (s *Store) Put(namespace, id, value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.key(namespace, id)
	s.values[key] = value
	return key
}

func (s *Store) Get(namespace, id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.values[s.key(namespace, id)]
	return value, ok
}

func (s *Store) Bump(namespace string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.versions[namespace]++
	return s.versions[namespace]
}

func (s *Store) Key(namespace, id string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.key(namespace, id)
}

func (s *Store) key(namespace, id string) string {
	return fmt.Sprintf("%s:v%d:%s", namespace, s.versions[namespace], id)
}
