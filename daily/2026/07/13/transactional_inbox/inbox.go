package transactionalinbox

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("message not found")
	ErrBusy     = errors.New("message lease still active")
	ErrDone     = errors.New("message already completed")
)

const (
	StatePending    = "pending"
	StateProcessing = "processing"
	StateDone       = "done"
)

type Message struct {
	ID         string
	Payload    string
	State      string
	Attempts   int
	LeaseUntil time.Time
	LastError  string
}

type Store struct {
	mu   sync.Mutex
	data map[string]Message
}

func NewStore() *Store {
	return &Store{data: make(map[string]Message)}
}

func (s *Store) Add(id string, payload string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = Message{
		ID:      id,
		Payload: payload,
		State:   StatePending,
	}
}

func (s *Store) Claim(id string, now time.Time, lease time.Duration) (Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.data[id]
	if !ok {
		return Message{}, ErrNotFound
	}
	if msg.State == StateDone {
		return Message{}, ErrDone
	}
	if msg.LeaseUntil.After(now) {
		return Message{}, ErrBusy
	}

	msg.State = StateProcessing
	msg.Attempts++
	msg.LeaseUntil = now.Add(lease)
	s.data[id] = msg
	return msg, nil
}

func (s *Store) Complete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	msg.State = StateDone
	msg.LeaseUntil = time.Time{}
	msg.LastError = ""
	s.data[id] = msg
	return nil
}

func (s *Store) Fail(id string, retryAt time.Time, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	if msg.State == StateDone {
		return ErrDone
	}
	msg.State = StatePending
	msg.LeaseUntil = retryAt
	msg.LastError = reason
	s.data[id] = msg
	return nil
}

func (s *Store) Get(id string) (Message, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.data[id]
	return msg, ok
}
