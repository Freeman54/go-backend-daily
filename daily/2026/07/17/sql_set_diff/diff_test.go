package sqlsetdiff

import (
	"reflect"
	"testing"
)

func TestBuildReturnsInsertAndDeleteSets(t *testing.T) {
	diff := Build([]int64{1, 2, 3}, []int64{2, 3, 4, 5})

	if !reflect.DeepEqual(diff.ToInsert, []int64{4, 5}) {
		t.Fatalf("unexpected inserts: %#v", diff.ToInsert)
	}
	if !reflect.DeepEqual(diff.ToDelete, []int64{1}) {
		t.Fatalf("unexpected deletes: %#v", diff.ToDelete)
	}
}

func TestBuildDeduplicatesInputAndSortsOutput(t *testing.T) {
	diff := Build([]int64{9, 9, 3, 1}, []int64{3, 3, 2, 10})

	if !reflect.DeepEqual(diff.ToInsert, []int64{2, 10}) {
		t.Fatalf("unexpected inserts: %#v", diff.ToInsert)
	}
	if !reflect.DeepEqual(diff.ToDelete, []int64{1, 9}) {
		t.Fatalf("unexpected deletes: %#v", diff.ToDelete)
	}
}
