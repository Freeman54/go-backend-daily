package fencingtoken

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLeaseHeld    = errors.New("lease is held by another owner")
	ErrLeaseExpired = errors.New("lease expired")
	ErrStaleToken   = errors.New("stale fencing token")
)

type Lease struct {
	Resource string
	Owner    string
	Token    uint64
	ExpireAt time.Time
}

type Manager struct {
	mu     sync.Mutex
	next   map[string]uint64
	active map[string]Lease
}

func NewManager() *Manager {
	return &Manager{
		next:   make(map[string]uint64),
		active: make(map[string]Lease),
	}
}

func (m *Manager) Grant(resource string, owner string, now time.Time, ttl time.Duration) (Lease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.active[resource]
	if ok && current.ExpireAt.After(now) && current.Owner != owner {
		return Lease{}, ErrLeaseHeld
	}

	nextToken := m.next[resource] + 1
	lease := Lease{
		Resource: resource,
		Owner:    owner,
		Token:    nextToken,
		ExpireAt: now.Add(ttl),
	}

	m.next[resource] = nextToken
	m.active[resource] = lease
	return lease, nil
}

func (m *Manager) ValidateWrite(resource string, token uint64, now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	current, ok := m.active[resource]
	if !ok {
		return ErrLeaseExpired
	}

	if !current.ExpireAt.After(now) {
		delete(m.active, resource)
		return ErrLeaseExpired
	}

	if current.Token != token {
		return ErrStaleToken
	}

	return nil
}
