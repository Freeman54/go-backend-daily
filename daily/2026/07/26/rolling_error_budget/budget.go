package rolling_error_budget

import (
	"errors"
	"sync"
	"time"
)

type observation struct {
	at      time.Time
	failure bool
}

// Budget 根据最近若干秒内的错误比例决定是否继续接收流量。
type Budget struct {
	mu           sync.Mutex
	window       time.Duration
	maxErrorRate float64
	observations []observation
}

func New(windowSeconds int, maxErrorRate float64) (*Budget, error) {
	if windowSeconds <= 0 || maxErrorRate < 0 || maxErrorRate > 1 {
		return nil, errors.New("invalid error budget configuration")
	}
	return &Budget{
		window:       time.Duration(windowSeconds) * time.Second,
		maxErrorRate: maxErrorRate,
	}, nil
}

func (b *Budget) Record(now time.Time, success bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.expire(now)
	b.observations = append(b.observations, observation{at: now, failure: !success})
}

func (b *Budget) Allow(now time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.expire(now)
	if len(b.observations) == 0 {
		return true
	}
	failures := 0
	for _, item := range b.observations {
		if item.failure {
			failures++
		}
	}
	return float64(failures)/float64(len(b.observations)) <= b.maxErrorRate
}

func (b *Budget) expire(now time.Time) {
	cutoff := now.Add(-b.window)
	first := 0
	for first < len(b.observations) && !b.observations[first].at.After(cutoff) {
		first++
	}
	b.observations = b.observations[first:]
}
