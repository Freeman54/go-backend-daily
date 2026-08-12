package responsebudget

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestWriterEnforcesBudget(t *testing.T) {
	recorder := httptest.NewRecorder()
	w := &Writer{ResponseWriter: recorder, Limit: 5}
	if n, err := w.Write([]byte("hello")); err != nil || n != 5 {
		t.Fatalf("first write = %d, %v", n, err)
	}
	if n, err := w.Write([]byte("!")); n != 0 || !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("overflow write = %d, %v", n, err)
	}
	if got := recorder.Body.String(); got != "hello" {
		t.Fatalf("body = %q", got)
	}
	if w.Written() != 5 {
		t.Fatalf("written = %d", w.Written())
	}
}

func TestWriterRejectsNegativeLimit(t *testing.T) {
	w := &Writer{ResponseWriter: httptest.NewRecorder(), Limit: -1}
	if _, err := w.Write(nil); !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("error = %v", err)
	}
}
