package messagesequenceguard

import (
	"errors"
	"fmt"
	"sync"
)

var ErrSequenceGap = errors.New("message sequence gap")

type Result uint8

const (
	Applied Result = iota + 1
	Duplicate
)

type Guard struct {
	mu   sync.Mutex
	last map[string]uint64
}

func New() *Guard {
	return &Guard{last: make(map[string]uint64)}
}

func (g *Guard) Accept(key string, sequence uint64) (Result, error) {
	if sequence == 0 {
		return 0, fmt.Errorf("sequence must be positive")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	last := g.last[key]
	if sequence <= last {
		return Duplicate, nil
	}
	if sequence != last+1 {
		return 0, fmt.Errorf("%w: got %d after %d", ErrSequenceGap, sequence, last)
	}
	g.last[key] = sequence
	return Applied, nil
}
