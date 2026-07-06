package savepointuow

import (
	"slices"
	"testing"
)

func TestPlannerNestedSuccess(t *testing.T) {
	var planner Planner

	sql, outer := planner.Enter()
	if !slices.Equal(sql, []string{"BEGIN"}) {
		t.Fatalf("outer enter = %v", sql)
	}

	sql, inner := planner.Enter()
	if !slices.Equal(sql, []string{"SAVEPOINT sp_1"}) {
		t.Fatalf("inner enter = %v", sql)
	}

	sql = planner.Leave(inner, true)
	if !slices.Equal(sql, []string{"RELEASE SAVEPOINT sp_1"}) {
		t.Fatalf("inner leave = %v", sql)
	}

	sql = planner.Leave(outer, true)
	if !slices.Equal(sql, []string{"COMMIT"}) {
		t.Fatalf("outer leave = %v", sql)
	}
}

func TestPlannerNestedRollbackDoesNotAbortOuterTransaction(t *testing.T) {
	var planner Planner

	_, outer := planner.Enter()
	_, inner := planner.Enter()

	sql := planner.Leave(inner, false)
	want := []string{"ROLLBACK TO SAVEPOINT sp_1", "RELEASE SAVEPOINT sp_1"}
	if !slices.Equal(sql, want) {
		t.Fatalf("inner rollback = %v, want %v", sql, want)
	}

	sql = planner.Leave(outer, true)
	if !slices.Equal(sql, []string{"COMMIT"}) {
		t.Fatalf("outer leave = %v", sql)
	}
}

func TestPlannerOuterRollback(t *testing.T) {
	var planner Planner

	_, outer := planner.Enter()
	sql := planner.Leave(outer, false)
	if !slices.Equal(sql, []string{"ROLLBACK"}) {
		t.Fatalf("outer leave = %v", sql)
	}
}
