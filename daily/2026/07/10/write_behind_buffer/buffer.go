package writebehindbuffer

import (
	"context"
	"sync"
	"time"
)

type Mutation struct {
	Key   string
	Delta int
}

type Buffer struct {
	mu        sync.Mutex
	threshold int
	pending   map[string]int
}

func New(threshold int) *Buffer {
	if threshold <= 0 {
		threshold = 1
	}
	return &Buffer{
		threshold: threshold,
		pending:   make(map[string]int),
	}
}

func (b *Buffer) Add(m Mutation) (flushNow bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending[m.Key] += m.Delta
	return len(b.pending) >= b.threshold
}

func (b *Buffer) Flush() map[string]int {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.pending) == 0 {
		return nil
	}

	flush := make(map[string]int, len(b.pending))
	for key, value := range b.pending {
		flush[key] = value
	}
	clear(b.pending)
	return flush
}

func (b *Buffer) Run(ctx context.Context, interval time.Duration, onFlush func(map[string]int)) {
	if interval <= 0 {
		interval = 50 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	flush := func() {
		if batch := b.Flush(); len(batch) > 0 {
			onFlush(batch)
		}
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case <-ticker.C:
			flush()
		}
	}
}
