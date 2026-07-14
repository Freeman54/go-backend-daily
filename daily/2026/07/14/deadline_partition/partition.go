package deadlinepartition

import (
	"context"
	"errors"
	"time"
)

var ErrNoDeadline = errors.New("context has no deadline")

type Step struct {
	Name   string
	Weight int
}

type Budget struct {
	Name    string
	Timeout time.Duration
}

func Partition(ctx context.Context, reserve time.Duration, steps []Step) ([]Budget, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, ErrNoDeadline
	}

	remaining := time.Until(deadline) - reserve
	if remaining < 0 {
		remaining = 0
	}

	totalWeight := 0
	for _, step := range steps {
		if step.Weight <= 0 {
			continue
		}
		totalWeight += step.Weight
	}

	budgets := make([]Budget, 0, len(steps))
	if totalWeight == 0 {
		for _, step := range steps {
			budgets = append(budgets, Budget{Name: step.Name})
		}
		return budgets, nil
	}

	allocated := time.Duration(0)
	for i, step := range steps {
		timeout := time.Duration(0)
		if step.Weight > 0 {
			if i == len(steps)-1 {
				timeout = remaining - allocated
			} else {
				timeout = remaining * time.Duration(step.Weight) / time.Duration(totalWeight)
				allocated += timeout
			}
		}
		if timeout < 0 {
			timeout = 0
		}
		budgets = append(budgets, Budget{Name: step.Name, Timeout: timeout})
	}
	return budgets, nil
}
