package optimisticlock

import (
	"errors"
	"sync"
)

var (
	ErrConflict = errors.New("version conflict")
	ErrNotFound = errors.New("record not found")
)

type Record struct {
	Value   string
	Version int64
}

type Store struct {
	mu   sync.Mutex
	data map[string]Record
}

func NewStore() *Store {
	return &Store{data: make(map[string]Record)}
}

func (s *Store) Create(key, value string) Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := Record{Value: value, Version: 1}
	s.data[key] = record
	return record
}

func (s *Store) Get(key string) (Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.data[key]
	if !ok {
		return Record{}, ErrNotFound
	}
	return record, nil
}

func (s *Store) Update(key string, expectedVersion int64, value string) (Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.data[key]
	if !ok {
		return Record{}, ErrNotFound
	}
	if current.Version != expectedVersion {
		return Record{}, ErrConflict
	}

	next := Record{
		Value:   value,
		Version: current.Version + 1,
	}
	s.data[key] = next
	return next, nil
}
