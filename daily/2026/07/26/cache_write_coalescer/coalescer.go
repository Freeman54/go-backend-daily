package cache_write_coalescer

import "sync"

type call struct {
	done chan struct{}
	err  error
}

// Coalescer 合并同一 key 上同时发生的等价写操作。
type Coalescer struct {
	mu       sync.Mutex
	inflight map[string]*call
}

func New() *Coalescer {
	return &Coalescer{inflight: make(map[string]*call)}
}

func (c *Coalescer) Do(key string, write func() error) error {
	c.mu.Lock()
	if existing, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		<-existing.done
		return existing.err
	}
	current := &call{done: make(chan struct{})}
	c.inflight[key] = current
	c.mu.Unlock()

	current.err = write()
	close(current.done)

	c.mu.Lock()
	delete(c.inflight, key)
	c.mu.Unlock()
	return current.err
}
