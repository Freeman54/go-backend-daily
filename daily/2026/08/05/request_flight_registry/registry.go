package request_flight_registry

import (
	"context"
	"sync"
)

type result[V any] struct {
	value V
	err   error
}

type flight[V any] struct {
	done    chan struct{}
	cancel  context.CancelFunc
	waiters int
	result  result[V]
}

type Registry[V any] struct {
	mu      sync.Mutex
	flights map[string]*flight[V]
}

func New[V any]() *Registry[V] {
	return &Registry[V]{flights: make(map[string]*flight[V])}
}

func (r *Registry[V]) Do(ctx context.Context, key string, work func(context.Context) (V, error)) (V, error) {
	r.mu.Lock()
	f := r.flights[key]
	if f == nil {
		workCtx, cancel := context.WithCancel(context.Background())
		f = &flight[V]{done: make(chan struct{}), cancel: cancel}
		r.flights[key] = f
		go func() {
			value, err := work(workCtx)
			r.mu.Lock()
			f.result = result[V]{value: value, err: err}
			delete(r.flights, key)
			close(f.done)
			r.mu.Unlock()
		}()
	}
	f.waiters++
	r.mu.Unlock()

	select {
	case <-f.done:
		return f.result.value, f.result.err
	case <-ctx.Done():
		r.mu.Lock()
		f.waiters--
		if f.waiters == 0 {
			f.cancel()
		}
		r.mu.Unlock()
		var zero V
		return zero, ctx.Err()
	}
}
