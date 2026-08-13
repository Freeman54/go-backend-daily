// Package cancelbarrier waits for concurrent work while propagating the first failure.
package cancelbarrier

import (
	"context"
	"sync"
)

// Run starts every function, cancels siblings after the first error, and waits for cleanup.
func Run(ctx context.Context, funcs ...func(context.Context) error) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	var wg sync.WaitGroup
	var once sync.Once
	var first error
	for _, fn := range funcs {
		fn := fn
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				once.Do(func() { first = err; cancel(err) })
			}
		}()
	}
	wg.Wait()
	if first != nil {
		return first
	}
	return context.Cause(ctx)
}
