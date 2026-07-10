package errorenvelope

import (
	"errors"
	"testing"
)

func TestBuildUsesTopLevelCodeAndMessage(t *testing.T) {
	root := errors.New("duplicate key")
	err := Wrap(CodeConflict, "email already exists", root)

	env := Build(err)
	if env.Code != CodeConflict || env.Message != "email already exists" {
		t.Fatalf("unexpected envelope: %+v", env)
	}
	if len(env.Causes) != 2 {
		t.Fatalf("expected 2 causes, got %#v", env.Causes)
	}
}

func TestBuildFallsBackToInternalForPlainErrors(t *testing.T) {
	env := Build(errors.New("dial tcp timeout"))
	if env.Code != CodeInternal || env.Message != "internal server error" {
		t.Fatalf("unexpected fallback envelope: %+v", env)
	}
	if len(env.Causes) != 1 || env.Causes[0] != "dial tcp timeout" {
		t.Fatalf("unexpected causes: %#v", env.Causes)
	}
}

func TestBuildDedupesRepeatedCauseStrings(t *testing.T) {
	leaf := errors.New("storage unavailable")
	err := Wrap(CodeUnauthorized, "token validation failed", Wrap(CodeInternal, "storage unavailable", leaf))

	env := Build(err)
	if len(env.Causes) != 2 {
		t.Fatalf("expected deduped causes, got %#v", env.Causes)
	}
}
