package tracelogenricher

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	traceIDKey   contextKey = "trace_id"
	spanIDKey    contextKey = "span_id"
	tenantKey    contextKey = "tenant"
	requestIDKey contextKey = "request_id"
)

func WithTraceID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, traceIDKey, value)
}

func WithSpanID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, spanIDKey, value)
}

func WithTenant(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, tenantKey, value)
}

func WithRequestID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, requestIDKey, value)
}

func AttrsFromContext(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0, 4)
	appendIfPresent := func(key contextKey) {
		if value, _ := ctx.Value(key).(string); value != "" {
			attrs = append(attrs, slog.String(string(key), value))
		}
	}

	appendIfPresent(traceIDKey)
	appendIfPresent(spanIDKey)
	appendIfPresent(tenantKey)
	appendIfPresent(requestIDKey)
	return attrs
}
