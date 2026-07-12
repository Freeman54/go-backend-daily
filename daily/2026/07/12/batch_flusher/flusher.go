package batchflusher

import (
	"sync"
	"time"
)

type FlushFunc[T any] func(items []T) error

type Flusher[T any] struct {
	mu       sync.Mutex
	maxBatch int
	maxWait  time.Duration
	flush    FlushFunc[T]
	batch    []T
	started  time.Time
}

func New[T any](maxBatch int, maxWait time.Duration, flush FlushFunc[T]) *Flusher[T] {
	return &Flusher[T]{
		maxBatch: maxBatch,
		maxWait:  maxWait,
		flush:    flush,
	}
}

func (f *Flusher[T]) Add(item T, now time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.batch) == 0 {
		f.started = now
	}
	f.batch = append(f.batch, item)
	if len(f.batch) < f.maxBatch {
		return nil
	}
	return f.flushLocked()
}

func (f *Flusher[T]) FlushDue(now time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.batch) == 0 || now.Sub(f.started) < f.maxWait {
		return nil
	}
	return f.flushLocked()
}

func (f *Flusher[T]) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.batch) == 0 {
		return nil
	}
	return f.flushLocked()
}

func (f *Flusher[T]) flushLocked() error {
	items := append([]T(nil), f.batch...)
	if err := f.flush(items); err != nil {
		return err
	}
	f.batch = nil
	f.started = time.Time{}
	return nil
}
