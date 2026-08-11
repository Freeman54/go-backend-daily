// Package generationlatch coordinates waiters across repeated state changes.
package generationlatch

import (
	"context"
	"sync"
)

// Latch lets callers wait for a generation strictly newer than the observed one.
type Latch struct {
	mu         sync.Mutex
	generation uint64
	changed    chan struct{}
}

func New() *Latch { return &Latch{changed: make(chan struct{})} }

func (l *Latch) Current() uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.generation
}

// Advance increments the generation and wakes all current waiters.
func (l *Latch) Advance() uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.generation++
	close(l.changed)
	l.changed = make(chan struct{})
	return l.generation
}

// Wait blocks until the generation is greater than observed or ctx is canceled.
func (l *Latch) Wait(ctx context.Context, observed uint64) (uint64, error) {
	for {
		l.mu.Lock()
		if l.generation > observed {
			generation := l.generation
			l.mu.Unlock()
			return generation, nil
		}
		changed := l.changed
		l.mu.Unlock()

		select {
		case <-changed:
		case <-ctx.Done():
			return observed, ctx.Err()
		}
	}
}
