package http_commit_guard

import (
	"net/http"
	"sync"
)

// Guard 确保多个竞争路径最多提交一次 HTTP 响应。
type Guard struct {
	mu        sync.Mutex
	writer    http.ResponseWriter
	committed bool
}

func New(writer http.ResponseWriter) *Guard {
	return &Guard{writer: writer}
}

func (g *Guard) Commit(status int, body []byte) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.committed {
		return false
	}
	g.committed = true
	g.writer.WriteHeader(status)
	_, _ = g.writer.Write(body)
	return true
}

func (g *Guard) Committed() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.committed
}
