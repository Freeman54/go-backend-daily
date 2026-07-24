package sqllockorder

import (
	"reflect"
	"testing"
)

func TestPlan_SortsAndDeduplicatesResourceIDs(t *testing.T) {
	got, err := Plan([]int64{9, 3, 9, 5})
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	want := []int64{3, 5, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Plan() = %v, want %v", got, want)
	}
}

func TestPlan_RejectsNonPositiveResourceID(t *testing.T) {
	if _, err := Plan([]int64{1, 0}); err == nil {
		t.Fatal("non-positive resource ID should return an error")
	}
}

func TestForUpdateQuery_BuildsStablePlaceholders(t *testing.T) {
	got, args, err := ForUpdateQuery("accounts", []int64{8, 2})
	if err != nil {
		t.Fatalf("ForUpdateQuery() error = %v", err)
	}
	if got != "SELECT id FROM accounts WHERE id IN ($1,$2) ORDER BY id FOR UPDATE" {
		t.Fatalf("query = %q", got)
	}
	if !reflect.DeepEqual(args, []int64{2, 8}) {
		t.Fatalf("args = %v, want [2 8]", args)
	}
}
