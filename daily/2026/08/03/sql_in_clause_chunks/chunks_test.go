package sql_in_clause_chunks

import (
	"reflect"
	"testing"
)

func TestPlan_SplitsValuesWithinParameterLimit(t *testing.T) {
	got, err := Plan([]int{1, 2, 3, 4, 5}, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Plan() = %#v, want %#v", got, want)
	}
}

func TestPlan_RejectsInvalidLimit(t *testing.T) {
	if _, err := Plan([]int{1}, 0); err == nil {
		t.Fatal("Plan() should reject zero limit")
	}
}
