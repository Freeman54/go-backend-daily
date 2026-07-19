package keyedsemaphore

import (
	"context"
	"fmt"
	"sync"
)

type Limiter struct {
	limit int

	mu      sync.Mutex
	buckets map[string]chan struct{}
}

func New(limit int) *Limiter {
	if limit <= 0 {
		panic("limit must be positive")
	}
	return &Limiter{
		limit:   limit,
		buckets: make(map[string]chan struct{}),
	}
}

func (l *Limiter) Acquire(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}

	bucket := l.bucket(key)
	select {
	case bucket <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Release(key string) error {
	l.mu.Lock()
	bucket, ok := l.buckets[key]
	l.mu.Unlock()
	if !ok {
		return fmt.Errorf("key %q not found", key)
	}

	select {
	case <-bucket:
		return nil
	default:
		return fmt.Errorf("key %q has no acquired slot", key)
	}
}

func (l *Limiter) bucket(key string) chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, ok := l.buckets[key]
	if ok {
		return bucket
	}

	bucket = make(chan struct{}, l.limit)
	l.buckets[key] = bucket
	return bucket
}
