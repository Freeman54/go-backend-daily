package deadline_cause_policy

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithTimeoutCause_ExposesBusinessCause(t *testing.T) {
	cause := errors.New("库存查询超时")
	ctx, cancel, err := NewTimeoutCause(context.Background(), time.Millisecond, cause)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	<-ctx.Done()
	if !errors.Is(context.Cause(ctx), cause) {
		t.Fatalf("Cause() = %v, want %v", context.Cause(ctx), cause)
	}
}

func TestWithTimeoutCause_RejectsInvalidTimeout(t *testing.T) {
	if _, _, err := NewTimeoutCause(context.Background(), 0, errors.New("x")); err == nil {
		t.Fatal("零超时时间应返回错误")
	}
	if _, _, err := NewTimeoutCause(context.Background(), time.Second, nil); err == nil {
		t.Fatal("nil cause 应返回错误")
	}
}
