package refreshaheadcache

import "time"

type Item struct {
	Value     string
	ExpiresAt time.Time
}

type Result struct {
	Value   string
	Hit     bool
	Stale   bool
	Refresh bool
}

type Cache struct {
	refreshAhead time.Duration
	item         Item
	hasItem      bool
	refreshing   bool
}

func New(refreshAhead time.Duration) *Cache {
	return &Cache{refreshAhead: refreshAhead}
}

func (c *Cache) Store(item Item) {
	c.item = item
	c.hasItem = true
	c.refreshing = false
}

func (c *Cache) Get(now time.Time) Result {
	if !c.hasItem {
		c.refreshing = true
		return Result{Refresh: true}
	}

	if !now.Before(c.item.ExpiresAt) {
		if !c.refreshing {
			c.refreshing = true
		}
		return Result{Refresh: true}
	}

	remaining := c.item.ExpiresAt.Sub(now)
	if remaining <= c.refreshAhead && !c.refreshing {
		c.refreshing = true
		return Result{
			Value:   c.item.Value,
			Hit:     true,
			Refresh: true,
		}
	}

	return Result{
		Value: c.item.Value,
		Hit:   true,
		Stale: remaining <= c.refreshAhead,
	}
}
