package httpproblemmapping

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}

type Problem struct {
	Status int
	Title  string
	Type   string
	Detail string
}

func Map(err error) Problem {
	if err == nil {
		return Problem{Status: http.StatusOK}
	}

	var validationErr ValidationError
	switch {
	case errors.As(err, &validationErr):
		return Problem{
			Status: http.StatusBadRequest,
			Title:  "invalid request",
			Type:   "https://example.dev/problems/validation",
			Detail: validationErr.Error(),
		}
	case errors.Is(err, ErrUnauthorized):
		return Problem{
			Status: http.StatusUnauthorized,
			Title:  "authentication required",
			Type:   "https://example.dev/problems/unauthorized",
			Detail: err.Error(),
		}
	case errors.Is(err, ErrNotFound):
		return Problem{
			Status: http.StatusNotFound,
			Title:  "resource not found",
			Type:   "https://example.dev/problems/not-found",
			Detail: err.Error(),
		}
	case errors.Is(err, ErrConflict):
		return Problem{
			Status: http.StatusConflict,
			Title:  "state conflict",
			Type:   "https://example.dev/problems/conflict",
			Detail: err.Error(),
		}
	default:
		return Problem{
			Status: http.StatusInternalServerError,
			Title:  "internal error",
			Type:   "https://example.dev/problems/internal",
			Detail: "please retry later",
		}
	}
}
