package httpproblemmapping

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestMapValidationError(t *testing.T) {
	problem := Map(ValidationError{
		Fields: map[string]string{"email": "required"},
	})

	if problem.Status != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", problem.Status, http.StatusBadRequest)
	}
	if problem.Title != "invalid request" {
		t.Fatalf("title = %q", problem.Title)
	}
}

func TestMapWrappedSentinel(t *testing.T) {
	problem := Map(fmt.Errorf("load account: %w", ErrNotFound))
	if problem.Status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", problem.Status, http.StatusNotFound)
	}
	if problem.Detail != "load account: not found" {
		t.Fatalf("detail = %q", problem.Detail)
	}
}

func TestMapUnknownErrorFallsBackToInternal(t *testing.T) {
	problem := Map(errors.New("sql: connection reset"))
	if problem.Status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", problem.Status, http.StatusInternalServerError)
	}
	if problem.Detail != "please retry later" {
		t.Fatalf("detail = %q", problem.Detail)
	}
}
