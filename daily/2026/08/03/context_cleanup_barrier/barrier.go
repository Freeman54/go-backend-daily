// Package context_cleanup_barrier waits for concurrent cleanup with a caller-owned deadline.
package context_cleanup_barrier

import (
	"context"
	"errors"
	"sync"
)

type Task func() error

func Wait(ctx context.Context, tasks ...Task) error {
	results := make(chan error, len(tasks))
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			results <- task()
		}(task)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	var joined error
	for {
		select {
		case <-ctx.Done():
			return errors.Join(joined, ctx.Err())
		case err, ok := <-results:
			if !ok {
				return joined
			}
			joined = errors.Join(joined, err)
		}
	}
}
