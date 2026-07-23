package boundedparallelmap

import (
	"context"
	"fmt"
)

type result[T any] struct {
	index int
	value T
	err   error
}

func Map[T, R any](ctx context.Context, input []T, limit int, fn func(context.Context, T) (R, error)) ([]R, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	output := make([]R, len(input))
	results := make(chan result[R], limit)
	next, running := 0, 0
	start := func(index int) {
		running++
		go func() {
			value, err := fn(ctx, input[index])
			results <- result[R]{index: index, value: value, err: err}
		}()
	}
	for next < len(input) && running < limit {
		start(next)
		next++
	}
	for running > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case current := <-results:
			running--
			if current.err != nil {
				return nil, current.err
			}
			output[current.index] = current.value
			if next < len(input) {
				start(next)
				next++
			}
		}
	}
	return output, nil
}
