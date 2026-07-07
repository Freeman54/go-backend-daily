package deadlineclamp

import (
	"context"
	"errors"
	"time"
)

func Clamp(parent context.Context, minBudget, maxBudget time.Duration) (context.Context, context.CancelFunc, time.Duration, error) {
	if minBudget <= 0 {
		return nil, nil, 0, errors.New("minBudget must be positive")
	}
	if maxBudget < minBudget {
		return nil, nil, 0, errors.New("maxBudget must be greater than or equal to minBudget")
	}

	if err := parent.Err(); err != nil {
		return nil, nil, 0, err
	}

	budget := maxBudget
	if deadline, ok := parent.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < minBudget {
			return nil, nil, 0, context.DeadlineExceeded
		}
		if remaining < budget {
			budget = remaining
		}
	}

	ctx, cancel := context.WithTimeout(parent, budget)
	return ctx, cancel, budget, nil
}
