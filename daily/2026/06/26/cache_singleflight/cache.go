package cachesingleflight

import (
	"context"
	"sync"
	"time"
)

type item[T any] struct {
	value     T
	expiresAt time.Time
}

type call[T any] struct {
	done  chan struct{}
	value T
	ttl   time.Duration
	err   error
}

type Cache[T any] struct {
	mu       sync.Mutex
	items    map[string]item[T]
	inflight map[string]*call[T]
	now      func() time.Time
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		items:    make(map[string]item[T]),
		inflight: make(map[string]*call[T]),
		now:      time.Now,
	}
}

func (c *Cache[T]) Get(ctx context.Context, key string, loader func(context.Context) (T, time.Duration, error)) (T, error) {
	if value, ok := c.lookup(key); ok {
		return value, nil
	}

	c.mu.Lock()
	if value, ok := c.lookupLocked(key); ok {
		c.mu.Unlock()
		return value, nil
	}
	if existing, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		select {
		case <-existing.done:
			return existing.value, existing.err
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}

	current := &call[T]{done: make(chan struct{})}
	c.inflight[key] = current
	c.mu.Unlock()

	current.value, current.ttl, current.err = loader(ctx)

	c.mu.Lock()
	delete(c.inflight, key)
	if current.err == nil {
		c.items[key] = item[T]{
			value:     current.value,
			expiresAt: c.now().Add(current.ttl),
		}
	}
	close(current.done)
	c.mu.Unlock()

	return current.value, current.err
}

func (c *Cache[T]) lookup(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lookupLocked(key)
}

func (c *Cache[T]) lookupLocked(key string) (T, bool) {
	item, ok := c.items[key]
	if !ok {
		var zero T
		return zero, false
	}
	if !c.now().Before(item.expiresAt) {
		delete(c.items, key)
		var zero T
		return zero, false
	}
	return item.value, true
}
