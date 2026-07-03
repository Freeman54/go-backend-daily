package negativecache

import (
	"sync"
	"time"
)

type State int

const (
	StateMiss State = iota
	StateHit
	StateNegativeHit
)

type entry struct {
	value    string
	state    State
	expireAt time.Time
}

type Cache struct {
	mu          sync.RWMutex
	now         func() time.Time
	positiveTTL time.Duration
	negativeTTL time.Duration
	entries     map[string]entry
}

func New(positiveTTL time.Duration, negativeTTL time.Duration) *Cache {
	return &Cache{
		now:         time.Now,
		positiveTTL: positiveTTL,
		negativeTTL: negativeTTL,
		entries:     make(map[string]entry),
	}
}

func (c *Cache) SetNow(now func() time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = now
}

func (c *Cache) SetFound(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{
		value:    value,
		state:    StateHit,
		expireAt: c.now().Add(c.positiveTTL),
	}
}

func (c *Cache) SetNotFound(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{
		state:    StateNegativeHit,
		expireAt: c.now().Add(c.negativeTTL),
	}
}

func (c *Cache) Get(key string) (string, State) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.entries[key]
	if !ok {
		return "", StateMiss
	}

	if !item.expireAt.After(c.now()) {
		delete(c.entries, key)
		return "", StateMiss
	}

	return item.value, item.state
}
