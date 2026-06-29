package gracefulshutdown

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Group struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewGroup() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *Group) Go(fn func(context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		_ = fn(g.ctx)
	}()
}

func (g *Group) Shutdown(timeout time.Duration) error {
	g.cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		g.wg.Wait()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
		return nil
	case <-timer.C:
		return errors.New("shutdown timed out")
	}
}
