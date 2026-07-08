package retrybudget

import (
	"sync"
	"time"
)

type snapshot struct {
	windowStart time.Time
	successes   int
	retries     int
}

type Budget struct {
	mu       sync.Mutex
	ratio    float64
	minRetry int
	window   time.Duration
	state    map[string]snapshot
}

func New(ratio float64, minRetry int, window time.Duration) *Budget {
	if ratio < 0 {
		ratio = 0
	}
	if minRetry < 0 {
		minRetry = 0
	}
	if window <= 0 {
		window = time.Minute
	}

	return &Budget{
		ratio:    ratio,
		minRetry: minRetry,
		window:   window,
		state:    make(map[string]snapshot),
	}
}

func (b *Budget) RecordSuccess(key string, now time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	current := b.rotate(key, now)
	current.successes++
	b.state[key] = current
}

func (b *Budget) AllowRetry(key string, now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	current := b.rotate(key, now)
	limit := b.minRetry + int(float64(current.successes)*b.ratio)
	if current.retries >= limit {
		b.state[key] = current
		return false
	}

	current.retries++
	b.state[key] = current
	return true
}

func (b *Budget) Snapshot(key string, now time.Time) (successes int, retries int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	current := b.rotate(key, now)
	b.state[key] = current
	return current.successes, current.retries
}

func (b *Budget) rotate(key string, now time.Time) snapshot {
	current, ok := b.state[key]
	if !ok || now.Sub(current.windowStart) >= b.window {
		return snapshot{windowStart: now}
	}
	return current
}
