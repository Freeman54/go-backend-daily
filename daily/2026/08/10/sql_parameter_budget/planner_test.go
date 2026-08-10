package parametbudget

import (
	"errors"
	"reflect"
	"testing"
)

func TestPlannerBatchSizes(t *testing.T) {
	p := Planner{MaxParameters: 10}
	got, err := p.BatchSizes(8, 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{3, 3, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BatchSizes() = %v, want %v", got, want)
	}

	empty, err := p.BatchSizes(0, 3, 1)
	if err != nil || len(empty) != 0 || empty == nil {
		t.Fatalf("empty BatchSizes() = %#v, %v", empty, err)
	}
}

func TestPlannerBatchSizesRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		planner Planner
		rows    int
		perRow  int
		fixed   int
	}{
		{Planner{}, 1, 1, 0},
		{Planner{10}, -1, 1, 0},
		{Planner{10}, 1, 0, 0},
		{Planner{10}, 1, 1, -1},
		{Planner{10}, 1, 1, 10},
		{Planner{10}, 1, 11, 0},
	}
	for _, tt := range tests {
		_, err := tt.planner.BatchSizes(tt.rows, tt.perRow, tt.fixed)
		if !errors.Is(err, ErrInvalidBudget) {
			t.Fatalf("BatchSizes() error = %v, want ErrInvalidBudget", err)
		}
	}
}
