package cache_negative_ttl

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value   V
	exists  bool
	expires time.Time
}

type Cache[V any] struct {
	mu          sync.Mutex
	items       map[string]entry[V]
	positiveTTL time.Duration
	negativeTTL time.Duration
	now         func() time.Time
}

func New[V any](positiveTTL, negativeTTL time.Duration, now func() time.Time) *Cache[V] {
	if positiveTTL <= 0 || negativeTTL <= 0 || now == nil {
		return nil
	}
	return &Cache[V]{items: make(map[string]entry[V]), positiveTTL: positiveTTL, negativeTTL: negativeTTL, now: now}
}

func (c *Cache[V]) Set(key string, value V, exists bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ttl := c.negativeTTL
	if exists {
		ttl = c.positiveTTL
	}
	c.items[key] = entry[V]{value: value, exists: exists, expires: c.now().Add(ttl)}
}

func (c *Cache[V]) Get(key string) (value V, exists, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if !found {
		return value, false, false
	}
	if !c.now().Before(item.expires) {
		delete(c.items, key)
		return value, false, false
	}
	return item.value, item.exists, true
}
