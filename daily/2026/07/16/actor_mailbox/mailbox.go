package actormailbox

import (
	"context"
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("mailbox queue full")

type Task func(context.Context)

type Mailbox struct {
	mu       sync.Mutex
	capacity int
	workers  map[string]*workerState
}

type workerState struct {
	running bool
	queue   []Task
}

func New(capacity int) *Mailbox {
	if capacity < 0 {
		capacity = 0
	}
	return &Mailbox{
		capacity: capacity,
		workers:  make(map[string]*workerState),
	}
}

func (m *Mailbox) Submit(ctx context.Context, key string, task Task) error {
	m.mu.Lock()
	state := m.workers[key]
	if state == nil {
		state = &workerState{}
		m.workers[key] = state
	}
	if state.running {
		if len(state.queue) >= m.capacity {
			m.mu.Unlock()
			return ErrQueueFull
		}
		state.queue = append(state.queue, task)
		m.mu.Unlock()
		return nil
	}
	state.running = true
	m.mu.Unlock()

	go m.run(ctx, key, task)
	return nil
}

func (m *Mailbox) run(ctx context.Context, key string, task Task) {
	for {
		task(ctx)

		m.mu.Lock()
		state := m.workers[key]
		if state == nil || len(state.queue) == 0 {
			delete(m.workers, key)
			m.mu.Unlock()
			return
		}
		task = state.queue[0]
		state.queue = state.queue[1:]
		m.mu.Unlock()
	}
}
