package cache_namespace_version

import "testing"

func TestNamespace_KeyUsesCurrentVersion(t *testing.T) {
	n := New("users")
	if got, want := n.Key("42"), "users:v1:42"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}
}

func TestNamespace_BumpInvalidatesPreviousKeys(t *testing.T) {
	n := New("users")
	old := n.Key("42")
	n.Bump()
	if got, want := n.Key("42"), "users:v2:42"; got != want || got == old {
		t.Fatalf("Key() = %q, want %q and different from %q", got, want, old)
	}
}

func TestNamespace_ConcurrentBumpsAreNotLost(t *testing.T) {
	n := New("users")
	done := make(chan struct{}, 20)
	for range 20 {
		go func() { n.Bump(); done <- struct{}{} }()
	}
	for range 20 {
		<-done
	}
	if got, want := n.Version(), uint64(21); got != want {
		t.Fatalf("Version() = %d, want %d", got, want)
	}
}
