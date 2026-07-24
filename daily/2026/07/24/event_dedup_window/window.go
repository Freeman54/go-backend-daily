package eventdedupwindow

import (
	"fmt"
	"sync"
	"time"
)

// Window 在进程内抑制指定时间范围内的重复事件。
type Window struct {
	mu      sync.Mutex
	ttl     time.Duration
	expires map[string]time.Time
}

func New(ttl time.Duration) (*Window, error) {
	if ttl <= 0 {
		return nil, fmt.Errorf("TTL must be positive")
	}
	return &Window{ttl: ttl, expires: make(map[string]time.Time)}, nil
}

func (w *Window) Accept(id string, now time.Time) bool {
	if id == "" {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	if expiry, exists := w.expires[id]; exists && now.Before(expiry) {
		return false
	}
	w.expires[id] = now.Add(w.ttl)
	return true
}
