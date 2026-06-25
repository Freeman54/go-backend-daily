package errortaxonomy

import (
	"errors"
	"fmt"
	"net/http"
)

type Kind string

const (
	KindInvalidArgument  Kind = "invalid_argument"
	KindUnauthenticated  Kind = "unauthenticated"
	KindPermissionDenied Kind = "permission_denied"
	KindNotFound         Kind = "not_found"
	KindConflict         Kind = "conflict"
	KindRateLimited      Kind = "rate_limited"
	KindInternal         Kind = "internal"
)

type AppError struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause == nil {
		return fmt.Sprintf("%s: %s", e.Kind, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Cause)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(kind Kind, message string) error {
	return &AppError{Kind: kind, Message: message}
}

func Wrap(kind Kind, message string, cause error) error {
	return &AppError{Kind: kind, Message: message, Cause: cause}
}

func KindOf(err error) Kind {
	if err == nil {
		return ""
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Kind
	}

	return KindInternal
}

func HTTPStatus(err error) int {
	switch KindOf(err) {
	case "":
		return http.StatusOK
	case KindInvalidArgument:
		return http.StatusBadRequest
	case KindUnauthenticated:
		return http.StatusUnauthorized
	case KindPermissionDenied:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

type PublicError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Public(err error) PublicError {
	kind := KindOf(err)
	if kind == "" {
		return PublicError{}
	}

	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Message != "" && kind != KindInternal {
		return PublicError{Code: string(kind), Message: appErr.Message}
	}

	return PublicError{Code: string(kind), Message: "服务暂时不可用"}
}
