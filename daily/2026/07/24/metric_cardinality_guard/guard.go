package metriccardinalityguard

import (
	"fmt"
	"sync"
)

// Guard 限制一个指标标签接受的不同值数量。
type Guard struct {
	mu       sync.Mutex
	limit    int
	fallback string
	known    map[string]struct{}
}

func New(limit int, fallback string) (*Guard, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive")
	}
	if fallback == "" {
		return nil, fmt.Errorf("fallback must not be empty")
	}
	return &Guard{limit: limit, fallback: fallback, known: make(map[string]struct{})}, nil
}

func (g *Guard) Normalize(value string) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.known[value]; exists {
		return value
	}
	if len(g.known) >= g.limit {
		return g.fallback
	}
	g.known[value] = struct{}{}
	return value
}
