package contextsnapshot

import (
	"context"
	"testing"
)

func TestCaptureCopiesWhitelistedValues(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-1")
	ctx = WithUserID(ctx, "user-9")
	ctx = WithTenant(ctx, "cn-hz")
	ctx = context.WithValue(ctx, contextKey("ignored"), "secret")

	got := Capture(ctx)

	if got.TraceID != "trace-1" || got.UserID != "user-9" || got.Tenant != "cn-hz" {
		t.Fatalf("unexpected snapshot: %+v", got)
	}
}

func TestAttachRestoresCapturedValues(t *testing.T) {
	snapshot := Snapshot{
		TraceID: "trace-1",
		UserID:  "user-9",
		Tenant:  "cn-hz",
	}

	ctx := snapshot.Attach(context.Background())

	if got := ctx.Value(traceIDKey); got != "trace-1" {
		t.Fatalf("unexpected trace id: %v", got)
	}
	if got := ctx.Value(userIDKey); got != "user-9" {
		t.Fatalf("unexpected user id: %v", got)
	}
	if got := ctx.Value(tenantKey); got != "cn-hz" {
		t.Fatalf("unexpected tenant: %v", got)
	}
}

func TestAttachSkipsEmptyValues(t *testing.T) {
	ctx := Snapshot{TraceID: "trace-1"}.Attach(context.Background())

	if got := ctx.Value(userIDKey); got != nil {
		t.Fatalf("expected empty value to stay absent, got %v", got)
	}
}
