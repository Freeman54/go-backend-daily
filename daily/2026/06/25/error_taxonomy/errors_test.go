package errortaxonomy

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestKindOfReadsWrappedAppError(t *testing.T) {
	err := fmt.Errorf("query user: %w", Wrap(KindNotFound, "用户不存在", sql.ErrNoRows))

	if got := KindOf(err); got != KindNotFound {
		t.Fatalf("expected %s, got %s", KindNotFound, got)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatal("expected wrapped error to keep original cause")
	}
}

func TestHTTPStatusMapsKnownKinds(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: http.StatusOK},
		{name: "invalid argument", err: New(KindInvalidArgument, "参数错误"), want: http.StatusBadRequest},
		{name: "unauthenticated", err: New(KindUnauthenticated, "请先登录"), want: http.StatusUnauthorized},
		{name: "permission denied", err: New(KindPermissionDenied, "无权限"), want: http.StatusForbidden},
		{name: "not found", err: New(KindNotFound, "资源不存在"), want: http.StatusNotFound},
		{name: "conflict", err: New(KindConflict, "状态冲突"), want: http.StatusConflict},
		{name: "rate limited", err: New(KindRateLimited, "请求过快"), want: http.StatusTooManyRequests},
		{name: "unknown", err: errors.New("redis unavailable"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTTPStatus(tt.err); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestPublicHidesInternalErrorMessage(t *testing.T) {
	err := Wrap(KindInternal, "mysql password leaked in detail", errors.New("connection refused"))

	got := Public(err)

	if got.Code != string(KindInternal) {
		t.Fatalf("expected code %s, got %s", KindInternal, got.Code)
	}
	if got.Message != "服务暂时不可用" {
		t.Fatalf("expected safe public message, got %q", got.Message)
	}
}

func TestPublicKeepsBusinessMessage(t *testing.T) {
	got := Public(New(KindInvalidArgument, "邮箱格式不正确"))

	if got.Code != string(KindInvalidArgument) {
		t.Fatalf("expected code %s, got %s", KindInvalidArgument, got.Code)
	}
	if got.Message != "邮箱格式不正确" {
		t.Fatalf("expected original business message, got %q", got.Message)
	}
}
