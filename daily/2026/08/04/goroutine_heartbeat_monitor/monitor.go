package goroutine_heartbeat_monitor

import (
	"sort"
	"sync"
	"time"
)

// Monitor records worker heartbeats and detects workers that have stopped making progress.
type Monitor struct {
	mu      sync.RWMutex
	timeout time.Duration
	beats   map[string]time.Time
}

func New(timeout time.Duration) *Monitor {
	if timeout <= 0 {
		panic("heartbeat timeout must be positive")
	}
	return &Monitor{timeout: timeout, beats: make(map[string]time.Time)}
}

func (m *Monitor) Beat(worker string, at time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.beats[worker] = at
}

func (m *Monitor) Remove(worker string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.beats, worker)
}

func (m *Monitor) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.beats)
}

func (m *Monitor) Stale(now time.Time) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stale := make([]string, 0)
	for worker, beat := range m.beats {
		if now.Sub(beat) > m.timeout {
			stale = append(stale, worker)
		}
	}
	sort.Strings(stale)
	return stale
}
