package gracefuldraintracker

import (
	"context"
	"fmt"
	"sync"
)

type Tracker struct {
	mu       sync.Mutex
	active   int
	draining bool
	drained  chan struct{}
}

func New() *Tracker { return &Tracker{drained: make(chan struct{})} }

func (t *Tracker) Begin() (func(), error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.draining {
		return nil, fmt.Errorf("tracker is draining")
	}
	t.active++
	var once sync.Once
	return func() { once.Do(t.finish) }, nil
}

func (t *Tracker) finish() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.active--
	if t.draining && t.active == 0 {
		close(t.drained)
	}
}

func (t *Tracker) Drain(ctx context.Context) error {
	t.mu.Lock()
	if !t.draining {
		t.draining = true
		if t.active == 0 {
			close(t.drained)
		}
	}
	drained := t.drained
	t.mu.Unlock()
	select {
	case <-drained:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
