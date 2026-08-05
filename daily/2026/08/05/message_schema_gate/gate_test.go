package message_schema_gate

import "testing"

func TestGateCheck_ClassifiesSchemaVersions(t *testing.T) {
	g, err := New(2, 4)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		version int
		want    Decision
	}{{1, Obsolete}, {2, Accept}, {4, Accept}, {5, Future}}
	for _, tt := range tests {
		if got := g.Check(tt.version); got != tt.want {
			t.Errorf("version %d: got %v, want %v", tt.version, got, tt.want)
		}
	}
}

func TestNew_RejectsInvalidRange(t *testing.T) {
	if _, err := New(0, 2); err == nil {
		t.Fatal("期望最小版本错误")
	}
	if _, err := New(3, 2); err == nil {
		t.Fatal("期望版本区间错误")
	}
}
