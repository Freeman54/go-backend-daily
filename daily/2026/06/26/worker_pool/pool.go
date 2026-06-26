package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	workers int
}

func New(workers int) Pool {
	if workers < 1 {
		workers = 1
	}
	return Pool{workers: workers}
}

func (p Pool) Run(ctx context.Context, jobs []int, fn func(context.Context, int) (int, error)) ([]int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		index int
		value int
		err   error
	}

	jobCh := make(chan int)
	resultCh := make(chan result, len(jobs))

	var workers sync.WaitGroup
	for range p.workers {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for index := range jobCh {
				value, err := fn(ctx, jobs[index])
				select {
				case resultCh <- result{index: index, value: value, err: err}:
				case <-ctx.Done():
				}
				if err != nil {
					cancel()
				}
			}
		}()
	}

	go func() {
		defer close(jobCh)
		for index := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- index:
			}
		}
	}()

	go func() {
		workers.Wait()
		close(resultCh)
	}()

	outputs := make([]int, len(jobs))
	for item := range resultCh {
		if item.err != nil {
			return nil, item.err
		}
		outputs[item.index] = item.value
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return outputs, nil
}
