package priorityadmission

import (
	"context"
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("priority admission queue full")

type Priority int

const (
	Low Priority = iota
	High
)

type request struct {
	done chan struct{}
}

type Limiter struct {
	mu              sync.Mutex
	totalInFlight   int
	highInFlight    int
	totalLimit      int
	reservedForHigh int
	queueLimit      int
	waitingHigh     []request
	waitingLow      []request
}

func New(totalLimit int, reservedForHigh int, queueLimit int) *Limiter {
	if totalLimit <= 0 {
		totalLimit = 1
	}
	if reservedForHigh < 0 {
		reservedForHigh = 0
	}
	if reservedForHigh > totalLimit {
		reservedForHigh = totalLimit
	}
	if queueLimit < 0 {
		queueLimit = 0
	}
	return &Limiter{
		totalLimit:      totalLimit,
		reservedForHigh: reservedForHigh,
		queueLimit:      queueLimit,
	}
}

func (l *Limiter) Acquire(ctx context.Context, p Priority) (func(), error) {
	l.mu.Lock()
	if l.canAdmitLocked(p) {
		l.takeLocked(p)
		l.mu.Unlock()
		return l.releaseFn(p), nil
	}
	if l.queueLimit > 0 && l.queueLenLocked() >= l.queueLimit {
		l.mu.Unlock()
		return nil, ErrQueueFull
	}

	req := request{done: make(chan struct{})}
	if p == High {
		l.waitingHigh = append(l.waitingHigh, req)
	} else {
		l.waitingLow = append(l.waitingLow, req)
	}
	l.mu.Unlock()

	select {
	case <-ctx.Done():
		l.mu.Lock()
		removed := l.removeWaitingLocked(p, req.done)
		l.mu.Unlock()
		if removed {
			return nil, ctx.Err()
		}
		return l.releaseFn(p), nil
	case <-req.done:
		return l.releaseFn(p), nil
	}
}

func (l *Limiter) releaseFn(p Priority) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			if l.totalInFlight > 0 {
				l.totalInFlight--
			}
			if p == High && l.highInFlight > 0 {
				l.highInFlight--
			}
			l.promoteWaitingLocked()
			l.mu.Unlock()
		})
	}
}

func (l *Limiter) canAdmitLocked(p Priority) bool {
	if l.totalInFlight >= l.totalLimit {
		return false
	}
	if p == High {
		return true
	}
	availableForLow := l.totalLimit - l.reservedForHigh
	if availableForLow < 0 {
		availableForLow = 0
	}
	lowInFlight := l.totalInFlight - l.highInFlight
	return lowInFlight < availableForLow
}

func (l *Limiter) takeLocked(p Priority) {
	l.totalInFlight++
	if p == High {
		l.highInFlight++
	}
}

func (l *Limiter) promoteWaitingLocked() {
	for {
		switch {
		case len(l.waitingHigh) > 0 && l.canAdmitLocked(High):
			req := l.waitingHigh[0]
			l.waitingHigh = l.waitingHigh[1:]
			l.takeLocked(High)
			close(req.done)
		case len(l.waitingLow) > 0 && l.canAdmitLocked(Low) && len(l.waitingHigh) == 0:
			req := l.waitingLow[0]
			l.waitingLow = l.waitingLow[1:]
			l.takeLocked(Low)
			close(req.done)
		default:
			return
		}
	}
}

func (l *Limiter) removeWaitingLocked(p Priority, done chan struct{}) bool {
	var queue *[]request
	if p == High {
		queue = &l.waitingHigh
	} else {
		queue = &l.waitingLow
	}
	for i, req := range *queue {
		if req.done == done {
			items := *queue
			*queue = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (l *Limiter) queueLenLocked() int {
	return len(l.waitingHigh) + len(l.waitingLow)
}
