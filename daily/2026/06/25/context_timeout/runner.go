package contexttimeout

import (
	"context"
	"errors"
	"sync"
)

var ErrInvalidWorkerCount = errors.New("worker count must be greater than zero")

type Job func(context.Context) error

type Result struct {
	Completed int
	Failed    int
	Canceled  bool
	Err       error
}

func RunWithWorkers(ctx context.Context, workers int, jobs []Job) Result {
	if workers <= 0 {
		return Result{Err: ErrInvalidWorkerCount}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobCh := make(chan Job)
	errCh := make(chan error, 1)
	doneCh := make(chan error, len(jobs))

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobCh {
				if err := job(ctx); err != nil {
					select {
					case errCh <- err:
						cancel()
					default:
					}
					doneCh <- err
					continue
				}
				doneCh <- nil
			}
		}()
	}

	go func() {
		defer close(jobCh)
		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- job:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	var result Result
	for err := range doneCh {
		if err != nil {
			result.Failed++
			if result.Err == nil {
				result.Err = err
			}
			continue
		}
		result.Completed++
	}

	select {
	case err := <-errCh:
		result.Err = err
	default:
	}

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			result.Canceled = true
			if result.Err == nil {
				result.Err = ctx.Err()
			}
		}
	default:
	}

	if result.Completed+result.Failed < len(jobs) {
		result.Canceled = true
		if result.Err == nil {
			result.Err = ctx.Err()
		}
	}

	return result
}
