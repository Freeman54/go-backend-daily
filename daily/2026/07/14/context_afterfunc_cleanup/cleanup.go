package contextafterfunccleanup

import (
	"context"
	"sync"
)

type Group struct {
	mu      sync.Mutex
	cleaned bool
	steps   []func()
}

func (g *Group) Add(step func()) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cleaned {
		step()
		return
	}
	g.steps = append(g.steps, step)
}

func (g *Group) Bind(ctx context.Context) func() bool {
	return context.AfterFunc(ctx, g.Cleanup)
}

func (g *Group) Cleanup() {
	g.mu.Lock()
	if g.cleaned {
		g.mu.Unlock()
		return
	}
	g.cleaned = true
	steps := make([]func(), len(g.steps))
	copy(steps, g.steps)
	g.steps = nil
	g.mu.Unlock()

	for i := len(steps) - 1; i >= 0; i-- {
		steps[i]()
	}
}
