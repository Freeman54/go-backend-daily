package batchloader

import (
	"context"
	"errors"
	"sync"
	"time"
)

type FetchFunc[K comparable, V any] func(context.Context, []K) (map[K]V, error)

type Loader[K comparable, V any] struct {
	mu        sync.Mutex
	wait      time.Duration
	maxBatch  int
	fetch     FetchFunc[K, V]
	pending   map[K][]chan result[V]
	flushChan chan struct{}
}

type result[V any] struct {
	value V
	err   error
}

func New[K comparable, V any](wait time.Duration, maxBatch int, fetch FetchFunc[K, V]) (*Loader[K, V], error) {
	if wait <= 0 {
		return nil, errors.New("wait must be positive")
	}
	if maxBatch <= 0 {
		return nil, errors.New("maxBatch must be positive")
	}
	if fetch == nil {
		return nil, errors.New("fetch must not be nil")
	}

	return &Loader[K, V]{
		wait:      wait,
		maxBatch:  maxBatch,
		fetch:     fetch,
		pending:   make(map[K][]chan result[V]),
		flushChan: make(chan struct{}, 1),
	}, nil
}

func (l *Loader[K, V]) Load(ctx context.Context, key K) (V, error) {
	reply := make(chan result[V], 1)
	shouldSchedule := false

	l.mu.Lock()
	if len(l.pending) == 0 {
		shouldSchedule = true
	}
	l.pending[key] = append(l.pending[key], reply)
	shouldFlushNow := len(l.pending) >= l.maxBatch
	l.mu.Unlock()

	if shouldSchedule {
		time.AfterFunc(l.wait, func() {
			l.signalFlush()
		})
	}
	if shouldFlushNow {
		l.signalFlush()
	}

	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	case outcome := <-reply:
		return outcome.value, outcome.err
	}
}

func (l *Loader[K, V]) Flush(ctx context.Context) error {
	for {
		keys, waiters := l.takeBatch()
		if len(keys) == 0 {
			return nil
		}
		if err := l.dispatch(ctx, keys, waiters); err != nil {
			return err
		}
	}
}

func (l *Loader[K, V]) signalFlush() {
	select {
	case l.flushChan <- struct{}{}:
	default:
	}
}

func (l *Loader[K, V]) Wait() <-chan struct{} {
	return l.flushChan
}

func (l *Loader[K, V]) takeBatch() ([]K, map[K][]chan result[V]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.pending) == 0 {
		return nil, nil
	}

	keys := make([]K, 0, l.maxBatch)
	waiters := make(map[K][]chan result[V], l.maxBatch)
	for key, list := range l.pending {
		keys = append(keys, key)
		waiters[key] = list
		delete(l.pending, key)
		if len(keys) == l.maxBatch {
			break
		}
	}
	return keys, waiters
}

func (l *Loader[K, V]) dispatch(ctx context.Context, keys []K, waiters map[K][]chan result[V]) error {
	values, err := l.fetch(ctx, keys)
	if err != nil {
		var zero V
		for _, key := range keys {
			for _, waiter := range waiters[key] {
				waiter <- result[V]{value: zero, err: err}
			}
		}
		return err
	}

	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			var zero V
			for _, waiter := range waiters[key] {
				waiter <- result[V]{value: zero, err: errors.New("missing key in fetch result")}
			}
			continue
		}
		for _, waiter := range waiters[key] {
			waiter <- result[V]{value: value}
		}
	}
	return nil
}
