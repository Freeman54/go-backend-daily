package cache_stale_if_error

import (
	"errors"
	"sync"
	"time"
)

type entry[T any] struct {
	value    T
	storedAt time.Time
}

// Cache 在刷新失败时允许返回仍处于兜底窗口内的旧值。
type Cache[T any] struct {
	mu       sync.RWMutex
	freshTTL time.Duration
	staleTTL time.Duration
	items    map[string]entry[T]
}

func New[T any](freshTTL, staleTTL time.Duration) (*Cache[T], error) {
	if freshTTL <= 0 || staleTTL < freshTTL {
		return nil, errors.New("invalid cache TTL")
	}
	return &Cache[T]{
		freshTTL: freshTTL,
		staleTTL: staleTTL,
		items:    make(map[string]entry[T]),
	}, nil
}

func (c *Cache[T]) Store(key string, value T, now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry[T]{value: value, storedAt: now}
}

func (c *Cache[T]) Get(key string, now time.Time, load func() (T, error)) (T, error) {
	c.mu.RLock()
	cached, found := c.items[key]
	c.mu.RUnlock()
	if found && now.Sub(cached.storedAt) <= c.freshTTL {
		return cached.value, nil
	}

	value, err := load()
	if err == nil {
		c.Store(key, value, now)
		return value, nil
	}
	if found && now.Sub(cached.storedAt) <= c.staleTTL {
		return cached.value, nil
	}
	var zero T
	return zero, err
}
