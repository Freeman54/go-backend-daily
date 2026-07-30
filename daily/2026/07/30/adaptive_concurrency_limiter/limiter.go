package adaptiveconcurrencylimiter

import (
	"errors"
	"sync"
	"time"
)

// Limiter 用一个采样窗口调整允许的并发数。
type Limiter struct {
	mu                        sync.Mutex
	limit, min, max, inFlight int
	target                    time.Duration
	samples, penalties        int
}

func New(initial, min, max int, target time.Duration) *Limiter {
	l, _ := NewChecked(initial, min, max, target)
	return l
}

func NewChecked(initial, min, max int, target time.Duration) (*Limiter, error) {
	if min < 1 || max < min || initial < min || initial > max || target <= 0 {
		return nil, errors.New("invalid limiter configuration")
	}
	return &Limiter{limit: initial, min: min, max: max, target: target}, nil
}

func (l *Limiter) TryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inFlight >= l.limit {
		return false
	}
	l.inFlight++
	return true
}

func (l *Limiter) Release(latency time.Duration, success bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.inFlight > 0 {
		l.inFlight--
	}
	l.samples++
	if !success || latency > l.target {
		l.penalties++
	}
	if l.samples < l.limit {
		return
	}
	if l.penalties == 0 && l.limit < l.max {
		l.limit++
	} else if l.penalties > 0 {
		l.limit -= l.penalties
		if l.limit < l.min {
			l.limit = l.min
		}
	}
	l.samples, l.penalties = 0, 0
}

func (l *Limiter) Limit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}
