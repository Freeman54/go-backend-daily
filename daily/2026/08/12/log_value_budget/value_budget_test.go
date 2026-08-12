package valuebudget

import (
	"testing"
	"unicode/utf8"
)

func TestTruncatePreservesUTF8(t *testing.T) {
	tests := []struct {
		value string
		max   int
		want  string
	}{
		{"ok", 2, "ok"},
		{"abcdef", 5, "ab…"},
		{"你好世界", 7, "你…"},
		{"abcdef", 2, ""},
		{"abcdef", 0, ""},
	}
	for _, tt := range tests {
		got := Truncate(tt.value, tt.max)
		if got != tt.want || !utf8.ValidString(got) || len(got) > tt.max {
			t.Fatalf("Truncate(%q, %d) = %q", tt.value, tt.max, got)
		}
	}
}

func TestFieldsReturnsDetachedMap(t *testing.T) {
	input := map[string]string{"error": "abcdefgh"}
	got := Fields(input, 6)
	if got["error"] != "abc…" {
		t.Fatalf("field = %q", got["error"])
	}
	got["error"] = "changed"
	if input["error"] != "abcdefgh" {
		t.Fatal("input map was mutated")
	}
}
