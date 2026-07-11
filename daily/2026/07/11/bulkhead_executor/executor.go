package bulkheadexecutor

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrLaneBusy = errors.New("bulkhead lane busy")

type Executor struct {
	mu    sync.RWMutex
	lanes map[string]chan struct{}
}

func New(config map[string]int) (*Executor, error) {
	if len(config) == 0 {
		return nil, errors.New("bulkhead config is empty")
	}

	lanes := make(map[string]chan struct{}, len(config))
	for lane, limit := range config {
		if limit <= 0 {
			return nil, fmt.Errorf("lane %q has non-positive limit", lane)
		}
		lanes[lane] = make(chan struct{}, limit)
	}

	return &Executor{lanes: lanes}, nil
}

func (e *Executor) Do(ctx context.Context, lane string, fn func(context.Context) error) error {
	token, ok := e.lane(lane)
	if !ok {
		return fmt.Errorf("lane %q not configured", lane)
	}

	select {
	case token <- struct{}{}:
		defer func() {
			<-token
		}()
	case <-ctx.Done():
		return errors.Join(ErrLaneBusy, ctx.Err())
	}

	return fn(ctx)
}

func (e *Executor) Occupancy(lane string) (used int, limit int, ok bool) {
	token, ok := e.lane(lane)
	if !ok {
		return 0, 0, false
	}
	return len(token), cap(token), true
}

func (e *Executor) lane(name string) (chan struct{}, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.lanes[name]
	return ch, ok
}
