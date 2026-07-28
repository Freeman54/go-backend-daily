package quorum_response_collector

import (
	"context"
	"errors"
)

var (
	ErrInvalidQuorum     = errors.New("invalid quorum")
	ErrQuorumUnavailable = errors.New("quorum unavailable")
)

type Task[T any] func(context.Context) (T, error)

type result[T any] struct {
	value T
	err   error
}

// Collect 并发执行任务，收集到 quorum 个成功结果后取消剩余任务。
func Collect[T any](parent context.Context, quorum int, tasks []Task[T]) ([]T, error) {
	if quorum <= 0 || quorum > len(tasks) {
		return nil, ErrInvalidQuorum
	}
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	results := make(chan result[T], len(tasks))
	for _, task := range tasks {
		go func(task Task[T]) {
			value, err := task(ctx)
			results <- result[T]{value: value, err: err}
		}(task)
	}

	values := make([]T, 0, quorum)
	failures := make([]error, 0, len(tasks)-quorum+1)
	for range tasks {
		outcome := <-results
		if outcome.err == nil {
			values = append(values, outcome.value)
			if len(values) == quorum {
				cancel()
				return values, nil
			}
			continue
		}
		failures = append(failures, outcome.err)
		if len(failures) > len(tasks)-quorum {
			return nil, errors.Join(append([]error{ErrQuorumUnavailable}, failures...)...)
		}
	}
	return nil, ErrQuorumUnavailable
}
