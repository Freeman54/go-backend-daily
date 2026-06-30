package tokenbucket

import (
	"sync"
	"time"
)

type Bucket struct {
	mu         sync.Mutex
	capacity   int
	tokens     float64
	refillRate float64
	lastRefill time.Time
}

func New(capacity int, refillEvery time.Duration, now time.Time) *Bucket {
	if capacity <= 0 {
		capacity = 1
	}
	rate := float64(capacity)
	if refillEvery > 0 {
		rate = float64(capacity) / refillEvery.Seconds()
	}

	return &Bucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: rate,
		lastRefill: now,
	}
}

func (b *Bucket) AllowN(now time.Time, n int) bool {
	if n <= 0 {
		return true
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill(now)
	if float64(n) > b.tokens {
		return false
	}
	b.tokens -= float64(n)
	return true
}

func (b *Bucket) Tokens(now time.Time) float64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill(now)
	return b.tokens
}

func (b *Bucket) refill(now time.Time) {
	if now.Before(b.lastRefill) {
		return
	}
	if b.refillRate <= 0 {
		b.tokens = float64(b.capacity)
		b.lastRefill = now
		return
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * b.refillRate
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.lastRefill = now
}
