package cachewritethrough

import (
	"errors"
	"sync"
)

var ErrStaleWrite = errors.New("stale write")

type Record struct {
	Value   string
	Version int64
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]Record),
	}
}

func (s *MemoryStore) Get(key string) (Record, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.data[key]
	return record, ok
}

func (s *MemoryStore) Upsert(key string, record Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.data[key]
	if ok && current.Version > record.Version {
		return
	}
	s.data[key] = record
}

func (s *MemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

type Store struct {
	repo  *MemoryStore
	cache *MemoryStore
}

func New(repo, cache *MemoryStore) *Store {
	return &Store{repo: repo, cache: cache}
}

func (s *Store) Save(key string, value string, version int64) error {
	current, ok := s.repo.Get(key)
	if ok && current.Version >= version {
		return ErrStaleWrite
	}

	record := Record{Value: value, Version: version}
	s.repo.Upsert(key, record)
	s.cache.Upsert(key, record)
	return nil
}

func (s *Store) Load(key string) (Record, bool) {
	if record, ok := s.cache.Get(key); ok {
		return record, true
	}

	record, ok := s.repo.Get(key)
	if !ok {
		return Record{}, false
	}

	s.cache.Upsert(key, record)
	return record, true
}
