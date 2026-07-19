package contextmerge

import (
	"context"
	"time"
)

func Merge(parents ...context.Context) (context.Context, context.CancelFunc) {
	var deadlineCancel context.CancelFunc = func() {}
	if deadline, ok := earliestDeadline(parents); ok {
		base, cancel := context.WithDeadline(context.Background(), deadline)
		deadlineCancel = cancel

		merged, cancelCause := context.WithCancelCause(base)
		for _, parent := range parents {
			if parent == nil {
				continue
			}
			go func(parent context.Context) {
				select {
				case <-parent.Done():
					cancelCause(context.Cause(parent))
				case <-merged.Done():
				}
			}(parent)
		}

		return merged, func() {
			deadlineCancel()
			cancelCause(context.Canceled)
		}
	}

	merged, cancelCause := context.WithCancelCause(context.Background())
	for _, parent := range parents {
		if parent == nil {
			continue
		}
		go func(parent context.Context) {
			select {
			case <-parent.Done():
				cancelCause(context.Cause(parent))
			case <-merged.Done():
			}
		}(parent)
	}

	return merged, func() {
		deadlineCancel()
		cancelCause(context.Canceled)
	}
}

func earliestDeadline(parents []context.Context) (time.Time, bool) {
	var earliest time.Time
	var ok bool
	for _, parent := range parents {
		if parent == nil {
			continue
		}
		deadline, hasDeadline := parent.Deadline()
		if !hasDeadline {
			continue
		}
		if !ok || deadline.Before(earliest) {
			earliest = deadline
			ok = true
		}
	}
	return earliest, ok
}
