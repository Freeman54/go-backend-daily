package querycostguard

import (
	"context"
	"fmt"
)

type Guard struct {
	tokens chan struct{}
	total  int
}

func New(total int) *Guard {
	if total <= 0 {
		panic("total must be positive")
	}

	tokens := make(chan struct{}, total)
	for i := 0; i < total; i++ {
		tokens <- struct{}{}
	}
	return &Guard{tokens: tokens, total: total}
}

func (g *Guard) Acquire(ctx context.Context, cost int) (func(), error) {
	if cost <= 0 {
		return nil, fmt.Errorf("cost must be positive")
	}
	if cost > g.total {
		return nil, fmt.Errorf("cost %d exceeds total budget %d", cost, g.total)
	}

	acquired := 0
	for acquired < cost {
		select {
		case <-g.tokens:
			acquired++
		case <-ctx.Done():
			for acquired > 0 {
				g.tokens <- struct{}{}
				acquired--
			}
			return nil, ctx.Err()
		}
	}

	released := false
	return func() {
		if released {
			return
		}
		released = true
		for i := 0; i < cost; i++ {
			g.tokens <- struct{}{}
		}
	}, nil
}
