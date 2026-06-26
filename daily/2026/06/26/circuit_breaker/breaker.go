package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker is open")

type State string

const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateHalfOpen State = "half-open"
)

type Config struct {
	FailureThreshold int
	HalfOpenAfter    time.Duration
	Now              func() time.Time
}

type Breaker struct {
	mu            sync.Mutex
	failureCount  int
	openedAt      time.Time
	state         State
	halfOpenProbe bool
	now           func() time.Time
	config        Config
}

func New(config Config) *Breaker {
	if config.FailureThreshold < 1 {
		config.FailureThreshold = 1
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	return &Breaker{
		state:  StateClosed,
		now:    config.Now,
		config: config,
	}
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refreshStateLocked()
	return b.state
}

func (b *Breaker) Execute(fn func() error) error {
	b.mu.Lock()
	b.refreshStateLocked()
	switch b.state {
	case StateOpen:
		b.mu.Unlock()
		return ErrOpen
	case StateHalfOpen:
		if b.halfOpenProbe {
			b.mu.Unlock()
			return ErrOpen
		}
		b.halfOpenProbe = true
	}
	b.mu.Unlock()

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()
	defer func() {
		if b.state == StateHalfOpen {
			b.halfOpenProbe = false
		}
	}()

	if err != nil {
		b.failureCount++
		if b.state == StateHalfOpen || b.failureCount >= b.config.FailureThreshold {
			b.state = StateOpen
			b.openedAt = b.now()
		}
		return err
	}

	b.failureCount = 0
	b.state = StateClosed
	b.openedAt = time.Time{}
	return nil
}

func (b *Breaker) refreshStateLocked() {
	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.config.HalfOpenAfter {
		b.state = StateHalfOpen
	}
}
