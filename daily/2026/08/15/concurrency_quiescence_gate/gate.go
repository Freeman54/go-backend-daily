package concurrencyquiescencegate

import (
	"errors"
	"sync"
)

var ErrClosed = errors.New("gate closed")

// Gate 让服务停止接收新任务，并等待已经进入的任务完成。
type Gate struct {
	mu     sync.Mutex
	closed bool
	active int
	done   chan struct{}
}

func New() *Gate { return &Gate{done: make(chan struct{})} }

// Enter 成功时返回幂等的 leave 函数。
func (g *Gate) Enter() (func(), error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return nil, ErrClosed
	}
	g.active++
	var once sync.Once
	return func() { once.Do(g.leave) }, nil
}

func (g *Gate) leave() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.active--
	if g.closed && g.active == 0 {
		close(g.done)
	}
}

// Close 停止接收新任务，返回在途任务清零后关闭的 channel。
func (g *Gate) Close() <-chan struct{} {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.closed {
		g.closed = true
		if g.active == 0 {
			close(g.done)
		}
	}
	return g.done
}

func (g *Gate) Active() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.active
}
