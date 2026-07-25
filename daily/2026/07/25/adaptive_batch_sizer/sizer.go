package adaptive_batch_sizer

import (
	"errors"
	"sync"
)

type Sizer struct {
	mu     sync.Mutex
	size   int
	step   int
	max    int
	target int
}

func New(initial, step, max, targetLatencyMS int) (*Sizer, error) {
	if initial <= 0 || step <= 0 || max < initial || targetLatencyMS <= 0 {
		return nil, errors.New("invalid sizing parameters")
	}
	return &Sizer{size: initial, step: step, max: max, target: targetLatencyMS}, nil
}

func (s *Sizer) Observe(latencyMS int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if latencyMS <= s.target {
		s.size = min(s.max, s.size+s.step)
	} else {
		s.size = max(1, s.size-s.step)
	}
	return s.size
}
