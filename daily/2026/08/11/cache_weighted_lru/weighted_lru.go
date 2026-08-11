// Package weightedlru implements a concurrency-safe cache bounded by entry weight.
package weightedlru

import (
	"container/list"
	"errors"
	"sync"
)

var ErrTooHeavy = errors.New("entry exceeds cache capacity")

type entry[K comparable, V any] struct {
	key    K
	value  V
	weight int64
}

// Cache evicts least-recently-used entries until the total weight fits.
type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int64
	weight   int64
	items    map[K]*list.Element
	order    *list.List
}

func New[K comparable, V any](capacity int64) (*Cache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	return &Cache[K, V]{capacity: capacity, items: make(map[K]*list.Element), order: list.New()}, nil
}

func (c *Cache[K, V]) Put(key K, value V, weight int64) error {
	if weight <= 0 {
		return errors.New("weight must be positive")
	}
	if weight > c.capacity {
		return ErrTooHeavy
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if old, ok := c.items[key]; ok {
		c.weight -= old.Value.(entry[K, V]).weight
		c.order.Remove(old)
	}
	elem := c.order.PushFront(entry[K, V]{key: key, value: value, weight: weight})
	c.items[key] = elem
	c.weight += weight
	for c.weight > c.capacity {
		victim := c.order.Back()
		item := victim.Value.(entry[K, V])
		delete(c.items, item.key)
		c.weight -= item.weight
		c.order.Remove(victim)
	}
	return nil
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.order.MoveToFront(elem)
	return elem.Value.(entry[K, V]).value, true
}

func (c *Cache[K, V]) Weight() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.weight
}
