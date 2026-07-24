package retrytokenbucket

import (
	"fmt"
	"sync"
	"time"
)

// Bucket 用固定间隔补充令牌，为额外重试建立独立预算。
type Bucket struct {
	mu       sync.Mutex
	capacity int
	tokens   int
	interval time.Duration
	last     time.Time
}

func New(capacity int, interval time.Duration, now time.Time) *Bucket {
	bucket, _ := NewChecked(capacity, interval, now)
	return bucket
}

func NewChecked(capacity int, interval time.Duration, now time.Time) (*Bucket, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be positive")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("refill interval must be positive")
	}
	return &Bucket{capacity: capacity, tokens: capacity, interval: interval, last: now}, nil
}

func (b *Bucket) Take(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if now.After(b.last) {
		refills := int(now.Sub(b.last) / b.interval)
		if refills > 0 {
			b.tokens = min(b.capacity, b.tokens+refills)
			b.last = b.last.Add(time.Duration(refills) * b.interval)
		}
	}
	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}
