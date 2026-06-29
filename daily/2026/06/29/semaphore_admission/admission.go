package semaphoreadmission

import "context"

type Gate struct {
	tokens chan struct{}
}

func NewGate(limit int) *Gate {
	if limit <= 0 {
		panic("limit must be positive")
	}

	return &Gate{
		tokens: make(chan struct{}, limit),
	}
}

func (g *Gate) Acquire(ctx context.Context) error {
	select {
	case g.tokens <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *Gate) Release() {
	select {
	case <-g.tokens:
	default:
		panic("release without acquire")
	}
}

func (g *Gate) InFlight() int {
	return len(g.tokens)
}

func (g *Gate) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := g.Acquire(ctx); err != nil {
		return err
	}
	defer g.Release()

	return fn(ctx)
}
