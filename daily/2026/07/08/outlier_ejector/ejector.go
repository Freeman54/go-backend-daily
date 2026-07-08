package outlierejector

import (
	"sync"
	"time"
)

type nodeState struct {
	consecutiveFailures int
	ejectedUntil        time.Time
}

type Ejector struct {
	mu               sync.Mutex
	failureThreshold int
	cooldown         time.Duration
	nodes            map[string]nodeState
}

func New(failureThreshold int, cooldown time.Duration) *Ejector {
	if failureThreshold <= 0 {
		failureThreshold = 1
	}
	if cooldown <= 0 {
		cooldown = 5 * time.Second
	}

	return &Ejector{
		failureThreshold: failureThreshold,
		cooldown:         cooldown,
		nodes:            make(map[string]nodeState),
	}
}

func (e *Ejector) Allow(node string, now time.Time) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.nodes[node]
	return !now.Before(state.ejectedUntil)
}

func (e *Ejector) Record(node string, success bool, now time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()

	state := e.nodes[node]
	if success {
		state.consecutiveFailures = 0
		state.ejectedUntil = time.Time{}
		e.nodes[node] = state
		return
	}

	state.consecutiveFailures++
	if state.consecutiveFailures >= e.failureThreshold {
		state.ejectedUntil = now.Add(e.cooldown)
		state.consecutiveFailures = 0
	}
	e.nodes[node] = state
}

func (e *Ejector) CooldownUntil(node string) time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.nodes[node].ejectedUntil
}
