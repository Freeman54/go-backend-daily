package context_first_success

import (
	"context"
	"errors"
)

var ErrNoTasks = errors.New("no tasks")

type Task[T any] func(context.Context) (T, error)

type result[T any] struct {
	value T
	err   error
}

// Run 并发执行任务，在第一个成功结果出现后取消其余任务。
func Run[T any](parent context.Context, tasks []Task[T]) (T, error) {
	var zero T
	if len(tasks) == 0 {
		return zero, ErrNoTasks
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

	failures := make([]error, 0, len(tasks))
	for range tasks {
		outcome := <-results
		if outcome.err == nil {
			cancel()
			return outcome.value, nil
		}
		failures = append(failures, outcome.err)
	}
	return zero, errors.Join(failures...)
}
