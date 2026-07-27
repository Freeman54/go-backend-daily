package http_commit_guard

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuard_CommitWritesOnlyFirstResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	guard := New(recorder)
	if !guard.Commit(http.StatusCreated, []byte("created")) {
		t.Fatal("first commit should succeed")
	}
	if guard.Commit(http.StatusInternalServerError, []byte("late error")) {
		t.Fatal("second commit should be ignored")
	}
	if recorder.Code != http.StatusCreated || recorder.Body.String() != "created" {
		t.Fatalf("response = %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestGuard_CommittedReportsState(t *testing.T) {
	guard := New(httptest.NewRecorder())
	if guard.Committed() {
		t.Fatal("new guard should not be committed")
	}
	guard.Commit(http.StatusNoContent, nil)
	if !guard.Committed() {
		t.Fatal("guard should report committed state")
	}
}
