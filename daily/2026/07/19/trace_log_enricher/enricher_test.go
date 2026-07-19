package tracelogenricher

import (
	"context"
	"testing"
)

func TestAttrsFromContextUsesWhitelistedKeys(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-1")
	ctx = WithSpanID(ctx, "span-1")
	ctx = WithTenant(ctx, "tenant-a")
	ctx = WithRequestID(ctx, "req-9")
	ctx = context.WithValue(ctx, contextKey("debug"), "ignored")

	attrs := AttrsFromContext(ctx)
	if len(attrs) != 4 {
		t.Fatalf("expected 4 attrs, got %d", len(attrs))
	}
	if attrs[0].Key != "trace_id" || attrs[0].Value.String() != "trace-1" {
		t.Fatalf("unexpected trace attr: %+v", attrs[0])
	}
}

func TestAttrsFromContextSkipsEmptyValues(t *testing.T) {
	ctx := WithTraceID(context.Background(), "")

	attrs := AttrsFromContext(ctx)
	if len(attrs) != 0 {
		t.Fatalf("expected no attrs, got %d", len(attrs))
	}
}
