package cachegenerationguard

import "sync"

type entry struct {
	generation uint64
	value      []byte
	present    bool
}

type Guard struct {
	mu      sync.RWMutex
	entries map[string]entry
}

func New() *Guard {
	return &Guard{entries: make(map[string]entry)}
}

func (g *Guard) Begin(key string) uint64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.entries[key].generation
}

func (g *Guard) Commit(key string, generation uint64, value []byte) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	current := g.entries[key]
	if current.generation != generation {
		return false
	}
	current.value = append([]byte(nil), value...)
	current.present = true
	g.entries[key] = current
	return true
}

func (g *Guard) Invalidate(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	current := g.entries[key]
	current.generation++
	current.value = nil
	current.present = false
	g.entries[key] = current
}

func (g *Guard) Get(key string) ([]byte, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	current := g.entries[key]
	if !current.present {
		return nil, false
	}
	return append([]byte(nil), current.value...), true
}
