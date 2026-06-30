package stalewhilerevalidate

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("cache entry not found")

type Entry struct {
	Value      string
	ExpiresAt  time.Time
	StaleUntil time.Time
}

type Result struct {
	Value         string
	Fresh         bool
	ShouldRefresh bool
}

type Cache struct {
	mu         sync.Mutex
	items      map[string]Entry
	refreshing map[string]bool
}

func New() *Cache {
	return &Cache{
		items:      make(map[string]Entry),
		refreshing: make(map[string]bool),
	}
}

func (c *Cache) Set(key string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry
	delete(c.refreshing, key)
}

func (c *Cache) Get(now time.Time, key string) (Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.items[key]
	if !ok {
		return Result{}, ErrNotFound
	}
	if now.Before(entry.ExpiresAt) || now.Equal(entry.ExpiresAt) {
		return Result{Value: entry.Value, Fresh: true}, nil
	}
	if now.Before(entry.StaleUntil) || now.Equal(entry.StaleUntil) {
		shouldRefresh := !c.refreshing[key]
		if shouldRefresh {
			c.refreshing[key] = true
		}
		return Result{Value: entry.Value, ShouldRefresh: shouldRefresh}, nil
	}
	return Result{}, ErrNotFound
}

func (c *Cache) Refresh(now time.Time, key string, ttl, staleWindow time.Duration, loader func() (string, error)) (string, error) {
	value, err := loader()
	if err != nil {
		c.mu.Lock()
		delete(c.refreshing, key)
		c.mu.Unlock()
		return "", err
	}

	c.Set(key, Entry{
		Value:      value,
		ExpiresAt:  now.Add(ttl),
		StaleUntil: now.Add(ttl + staleWindow),
	})
	return value, nil
}
