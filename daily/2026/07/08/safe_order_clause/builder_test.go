package safeorderclause

import "testing"

func TestBuildAscendingOrderClause(t *testing.T) {
	builder := New(map[string]string{
		"created_at": "orders.created_at",
		"amount":     "orders.amount",
	})

	clause, err := builder.Build("created_at", false)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if clause != "orders.created_at ASC, id ASC" {
		t.Fatalf("unexpected clause: %s", clause)
	}
}

func TestBuildDescendingOrderClause(t *testing.T) {
	builder := New(map[string]string{
		"amount": "orders.amount",
	})

	clause, err := builder.Build("amount", true)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if clause != "orders.amount DESC, id ASC" {
		t.Fatalf("unexpected clause: %s", clause)
	}
}

func TestBuildRejectsUnsupportedField(t *testing.T) {
	builder := New(map[string]string{
		"amount": "orders.amount",
	})

	if _, err := builder.Build("drop table orders", true); err == nil {
		t.Fatal("expected unsupported field error")
	}
}
