package sqlsavepointplan

import "testing"

func TestPlan_ProducesQuotedSavepointStatements(t *testing.T) {
	p, err := New("import_batch_1")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := p.Begin(), `SAVEPOINT "import_batch_1"`; got != want {
		t.Fatalf("Begin() = %q, want %q", got, want)
	}
	if got, want := p.Rollback(), `ROLLBACK TO SAVEPOINT "import_batch_1"`; got != want {
		t.Fatalf("Rollback() = %q, want %q", got, want)
	}
	if got, want := p.Release(), `RELEASE SAVEPOINT "import_batch_1"`; got != want {
		t.Fatalf("Release() = %q, want %q", got, want)
	}
}

func TestNew_RejectsUnsafeIdentifier(t *testing.T) {
	for _, name := range []string{"", "a-b", `x"; DROP TABLE users;--`, "包含中文"} {
		if _, err := New(name); err == nil {
			t.Fatalf("New(%q) 应拒绝不安全标识符", name)
		}
	}
}
