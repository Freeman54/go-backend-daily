package messagecorrelationtracker

import (
	"errors"
	"sync"
)

var (
	ErrDuplicate = errors.New("duplicate correlation id")
	ErrUnknown   = errors.New("unknown correlation id")
)

type Result struct {
	Payload []byte
	Err     error
}

type Tracker struct {
	mu      sync.Mutex
	pending map[string]chan Result
}

func New() *Tracker { return &Tracker{pending: make(map[string]chan Result)} }

// Register 必须在发布请求消息前调用，避免快速响应先于注册到达。
func (t *Tracker) Register(id string) (<-chan Result, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if id == "" {
		return nil, ErrUnknown
	}
	if _, exists := t.pending[id]; exists {
		return nil, ErrDuplicate
	}
	ch := make(chan Result, 1)
	t.pending[id] = ch
	return ch, nil
}

func (t *Tracker) Resolve(id string, result Result) error {
	t.mu.Lock()
	ch, ok := t.pending[id]
	if ok {
		delete(t.pending, id)
	}
	t.mu.Unlock()
	if !ok {
		return ErrUnknown
	}
	ch <- result
	close(ch)
	return nil
}

// Cancel 在调用方超时或取消时释放登记项。
func (t *Tracker) Cancel(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.pending[id]; !ok {
		return false
	}
	delete(t.pending, id)
	return true
}

func (t *Tracker) Pending() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.pending)
}
