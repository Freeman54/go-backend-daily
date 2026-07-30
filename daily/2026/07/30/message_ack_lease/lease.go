package messageacklease

import (
	"errors"
	"sync"
	"time"
)

type Lease struct {
	mu           sync.Mutex
	deadline     time.Time
	safetyMargin time.Duration
	done         bool
}

func New(now time.Time, duration, safetyMargin time.Duration) (*Lease, error) {
	if duration <= 0 || safetyMargin <= 0 || safetyMargin >= duration {
		return nil, errors.New("invalid lease duration or safety margin")
	}
	return &Lease{deadline: now.Add(duration), safetyMargin: safetyMargin}, nil
}

func (l *Lease) ShouldExtend(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return !l.done && !now.Before(l.deadline.Add(-l.safetyMargin))
}

func (l *Lease) Extended(now time.Time, duration time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.deadline = now.Add(duration)
}

func (l *Lease) Ack() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.done = true
}
