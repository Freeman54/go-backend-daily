package cachegenerationguard

import "testing"

func TestGuardRejectsStaleFillAfterInvalidation(t *testing.T) {
	g := New()
	generation := g.Begin("user:42")
	g.Invalidate("user:42")
	if g.Commit("user:42", generation, []byte("stale")) {
		t.Fatal("stale fill was committed")
	}
	if _, ok := g.Get("user:42"); ok {
		t.Fatal("stale value became visible")
	}
}

func TestGuardCommitsCurrentFillAndCopiesBytes(t *testing.T) {
	g := New()
	value := []byte("fresh")
	if !g.Commit("user:42", g.Begin("user:42"), value) {
		t.Fatal("current fill was rejected")
	}
	value[0] = 'X'
	got, ok := g.Get("user:42")
	if !ok || string(got) != "fresh" {
		t.Fatalf("Get() = %q, %v; want fresh, true", got, ok)
	}
	got[0] = 'Y'
	again, _ := g.Get("user:42")
	if string(again) != "fresh" {
		t.Fatalf("stored bytes mutated through caller: %q", again)
	}
}

func TestInvalidateRemovesCommittedValue(t *testing.T) {
	g := New()
	g.Commit("key", g.Begin("key"), []byte("value"))
	g.Invalidate("key")
	if _, ok := g.Get("key"); ok {
		t.Fatal("invalidated value still exists")
	}
}
