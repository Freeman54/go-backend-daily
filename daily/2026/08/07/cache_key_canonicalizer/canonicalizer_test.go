package cachekeycanonicalizer

import "testing"

func TestBuild_SortsAndEscapesParts(t *testing.T) {
	t.Parallel()
	got, err := Build("product", map[string]string{"locale": "zh CN", "id": "42"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if want := "product?id=42&locale=zh+CN"; got != want {
		t.Errorf("Build() = %q, want %q", got, want)
	}
}

func TestBuild_RejectsEmptyNamespace(t *testing.T) {
	t.Parallel()
	if _, err := Build("", map[string]string{"id": "42"}); err == nil {
		t.Fatal("Build() error = nil, want namespace error")
	}
}

func TestBuild_RejectsEmptyPartName(t *testing.T) {
	t.Parallel()
	if _, err := Build("product", map[string]string{"": "42"}); err == nil {
		t.Fatal("Build() error = nil, want part name error")
	}
}
