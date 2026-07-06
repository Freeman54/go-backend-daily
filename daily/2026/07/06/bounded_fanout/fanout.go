package boundedfanout

import (
	"context"
	"errors"
	"sync"
)

type Task func(context.Context) (string, error)

func Run(ctx context.Context, maxParallel int, tasks []Task) ([]string, error) {
	if maxParallel <= 0 {
		return nil, errors.New("maxParallel must be positive")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]string, len(tasks))
	sem := make(chan struct{}, maxParallel)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i, task := range tasks {
		if err := ctx.Err(); err != nil {
			break
		}

		wg.Add(1)
		go func(idx int, current Task) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() {
				<-sem
			}()

			value, err := current(ctx)
			if err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
				return
			}
			results[idx] = value
		}(i, task)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return nil, err
	}
	return results, nil
}
