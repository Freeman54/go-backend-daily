package methodpolicy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPolicyMiddleware(t *testing.T) {
	policy, err := New("post", "GET", "GET")
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })

	allowed := httptest.NewRecorder()
	policy.Middleware(next).ServeHTTP(allowed, httptest.NewRequest(http.MethodGet, "/", nil))
	if allowed.Code != http.StatusNoContent {
		t.Fatalf("allowed status = %d", allowed.Code)
	}

	rejected := httptest.NewRecorder()
	policy.Middleware(next).ServeHTTP(rejected, httptest.NewRequest(http.MethodDelete, "/", nil))
	if rejected.Code != http.StatusMethodNotAllowed || rejected.Header().Get("Allow") != "GET, POST" {
		t.Fatalf("rejected status/header = %d/%q", rejected.Code, rejected.Header().Get("Allow"))
	}
}

func TestPolicyValidation(t *testing.T) {
	if _, err := New(); err == nil {
		t.Fatal("New() expected error")
	}
	if _, err := New("GET /admin"); err == nil {
		t.Fatal("New(invalid) expected error")
	}
	policy, _ := New("GET")
	_, err := policy.Check("PATCH")
	if !errors.Is(err, ErrMethodNotAllowed) {
		t.Fatalf("error = %v", err)
	}
}
