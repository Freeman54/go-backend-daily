package timeoutbudget

import (
	"context"
	"errors"
	"time"
)

var ErrNoDeadline = errors.New("context has no deadline")

type Step struct {
	Name     string
	Weight   int
	MinShare time.Duration
}

type Allocation struct {
	Name    string
	Timeout time.Duration
}

func Split(ctx context.Context, now time.Time, steps []Step) ([]Allocation, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, ErrNoDeadline
	}

	remaining := deadline.Sub(now)
	if remaining <= 0 {
		return nil, context.DeadlineExceeded
	}

	totalMin := time.Duration(0)
	totalWeight := 0
	for _, step := range steps {
		totalMin += step.MinShare
		if step.Weight > 0 {
			totalWeight += step.Weight
		}
	}

	if totalMin > remaining {
		return nil, context.DeadlineExceeded
	}

	extra := remaining - totalMin
	out := make([]Allocation, 0, len(steps))
	distributed := time.Duration(0)
	for i, step := range steps {
		timeout := step.MinShare
		if totalWeight > 0 && step.Weight > 0 {
			share := time.Duration(int64(extra) * int64(step.Weight) / int64(totalWeight))
			timeout += share
			distributed += share
		}

		if i == len(steps)-1 {
			timeout += extra - distributed
		}

		out = append(out, Allocation{
			Name:    step.Name,
			Timeout: timeout,
		})
	}
	return out, nil
}
