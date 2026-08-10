// Package leaktracker tracks the lifetime of application-owned background tasks.
package leaktracker

import (
	"sort"
	"sync"
	"time"
)

type Task struct {
	ID      uint64
	Name    string
	Started time.Time
}

type Tracker struct {
	mu     sync.Mutex
	nextID uint64
	tasks  map[uint64]Task
	now    func() time.Time
}

func New() *Tracker {
	return &Tracker{tasks: make(map[uint64]Task), now: time.Now}
}

// Start registers a task and returns an idempotent completion function.
func (t *Tracker) Start(name string) func() {
	t.mu.Lock()
	t.nextID++
	id := t.nextID
	t.tasks[id] = Task{ID: id, Name: name, Started: t.now()}
	t.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			t.mu.Lock()
			delete(t.tasks, id)
			t.mu.Unlock()
		})
	}
}

// OlderThan returns tasks at least maxAge old, oldest first.
func (t *Tracker) OlderThan(now time.Time, maxAge time.Duration) []Task {
	t.mu.Lock()
	defer t.mu.Unlock()
	result := make([]Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		if maxAge <= 0 || now.Sub(task.Started) >= maxAge {
			result = append(result, task)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Started.Equal(result[j].Started) {
			return result[i].ID < result[j].ID
		}
		return result[i].Started.Before(result[j].Started)
	})
	return result
}
