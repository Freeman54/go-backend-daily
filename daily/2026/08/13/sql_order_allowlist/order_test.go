package orderallowlist

import "testing"

func TestBuild(t *testing.T) {
	allowed := map[string]string{"createdAt": "created_at", "id": "id"}
	got, err := Build([]Field{{"createdAt", Desc}, {"id", ""}}, allowed)
	if err != nil || got != "ORDER BY created_at DESC, id ASC" {
		t.Fatalf("got %q, %v", got, err)
	}
}
func TestBuildRejectsInput(t *testing.T) {
	allowed := map[string]string{"id": "id"}
	for _, fields := range [][]Field{{{"id; DROP TABLE users", Asc}}, {{"id", Direction("SIDEWAYS")}}} {
		if _, err := Build(fields, allowed); err == nil {
			t.Fatal("expected error")
		}
	}
	if got, err := Build(nil, allowed); err != nil || got != "" {
		t.Fatalf("got %q, %v", got, err)
	}
}
