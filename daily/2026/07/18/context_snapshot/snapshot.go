package contextsnapshot

import "context"

type contextKey string

const (
	traceIDKey contextKey = "trace_id"
	userIDKey  contextKey = "user_id"
	tenantKey  contextKey = "tenant"
)

type Snapshot struct {
	TraceID string
	UserID  string
	Tenant  string
}

func WithTraceID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, traceIDKey, value)
}

func WithUserID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, userIDKey, value)
}

func WithTenant(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, tenantKey, value)
}

func Capture(ctx context.Context) Snapshot {
	return Snapshot{
		TraceID: stringValue(ctx, traceIDKey),
		UserID:  stringValue(ctx, userIDKey),
		Tenant:  stringValue(ctx, tenantKey),
	}
}

func (s Snapshot) Attach(ctx context.Context) context.Context {
	if s.TraceID != "" {
		ctx = WithTraceID(ctx, s.TraceID)
	}
	if s.UserID != "" {
		ctx = WithUserID(ctx, s.UserID)
	}
	if s.Tenant != "" {
		ctx = WithTenant(ctx, s.Tenant)
	}
	return ctx
}

func stringValue(ctx context.Context, key contextKey) string {
	value, _ := ctx.Value(key).(string)
	return value
}
