// Package entrycostbudget decides whether a cache entry is worth admitting.
package entrycostbudget

import "sync"

// Budget tracks the total cost of admitted cache entries.
type Budget struct {
	mu    sync.Mutex
	limit int64
	used  int64
}

func New(limit int64) *Budget { return &Budget{limit: limit} }

// Reserve atomically admits a positive cost when enough capacity remains.
func (b *Budget) Reserve(cost int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if cost <= 0 || b.limit < 0 || cost > b.limit-b.used {
		return false
	}
	b.used += cost
	return true
}

// Release returns previously reserved cost and clamps misuse at zero.
func (b *Budget) Release(cost int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if cost >= b.used {
		b.used = 0
		return
	}
	if cost > 0 {
		b.used -= cost
	}
}

func (b *Budget) Used() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.used
}
