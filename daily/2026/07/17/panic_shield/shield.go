package panicshield

import (
	"context"
	"fmt"
	"runtime/debug"
)

type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic recovered: %v", e.Value)
}

func Execute(ctx context.Context, fn func(context.Context) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Value: r,
				Stack: debug.Stack(),
			}
		}
	}()

	return fn(ctx)
}
