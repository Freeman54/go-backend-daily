package readiness_hysteresis

import (
	"errors"
	"sync"
)

// Gate 用不同的失败和恢复阈值抑制就绪状态反复抖动。
type Gate struct {
	mu                sync.Mutex
	ready             bool
	failureThreshold  int
	recoveryThreshold int
	failures          int
	successes         int
}

func New(failureThreshold, recoveryThreshold int) (*Gate, error) {
	if failureThreshold <= 0 || recoveryThreshold <= 0 {
		return nil, errors.New("thresholds must be positive")
	}
	return &Gate{
		ready:             true,
		failureThreshold:  failureThreshold,
		recoveryThreshold: recoveryThreshold,
	}, nil
}

func (g *Gate) Observe(success bool) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if success {
		g.failures = 0
		g.successes++
		if !g.ready && g.successes >= g.recoveryThreshold {
			g.ready = true
		}
		return g.ready
	}

	g.successes = 0
	g.failures++
	if g.ready && g.failures >= g.failureThreshold {
		g.ready = false
	}
	return g.ready
}
