package sagacompensation

import (
	"context"
	"errors"
	"fmt"
)

type Step struct {
	Name       string
	Do         func(context.Context) error
	Compensate func(context.Context) error
}

func Execute(ctx context.Context, steps []Step) error {
	completed := make([]Step, 0, len(steps))

	for _, step := range steps {
		if err := step.Do(ctx); err != nil {
			return compensate(ctx, completed, fmt.Errorf("step %s failed: %w", step.Name, err))
		}
		completed = append(completed, step)
	}

	return nil
}

func compensate(ctx context.Context, completed []Step, original error) error {
	errs := []error{original}
	for i := len(completed) - 1; i >= 0; i-- {
		step := completed[i]
		if step.Compensate == nil {
			continue
		}
		if err := step.Compensate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("compensate %s failed: %w", step.Name, err))
		}
	}
	return errors.Join(errs...)
}
