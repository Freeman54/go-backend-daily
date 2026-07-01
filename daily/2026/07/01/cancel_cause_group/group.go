package cancelcausegroup

import (
	"context"
	"sync"
)

type Task func(context.Context) error

func Run(ctx context.Context, tasks ...Task) error {
	groupCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	var once sync.Once
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(run Task) {
			defer wg.Done()
			if err := run(groupCtx); err != nil {
				once.Do(func() {
					cancel(err)
				})
			}
		}(task)
	}

	wg.Wait()
	return context.Cause(groupCtx)
}
