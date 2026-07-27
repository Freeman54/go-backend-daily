package sql_filter_builder

import (
	"reflect"
	"testing"
)

func TestBuild_CreatesStableParameterizedPredicates(t *testing.T) {
	query, args, err := Build("SELECT id FROM orders", []Filter{
		{Column: "tenant_id", Value: 7},
		{Column: "status", Value: "paid"},
	}, DialectPostgres)
	if err != nil {
		t.Fatal(err)
	}
	if query != "SELECT id FROM orders WHERE tenant_id = $1 AND status = $2" {
		t.Fatalf("Build() query = %q", query)
	}
	if !reflect.DeepEqual(args, []any{7, "paid"}) {
		t.Fatalf("Build() args = %#v", args)
	}
}

func TestBuild_UsesQuestionMarksForMySQL(t *testing.T) {
	query, args, err := Build("SELECT id FROM orders", []Filter{
		{Column: "status", Value: "paid"},
	}, DialectMySQL)
	if err != nil {
		t.Fatal(err)
	}
	if query != "SELECT id FROM orders WHERE status = ?" || len(args) != 1 {
		t.Fatalf("Build() = (%q, %#v)", query, args)
	}
}

func TestBuild_RejectsUnsafeColumnName(t *testing.T) {
	_, _, err := Build("SELECT id FROM orders", []Filter{
		{Column: "status OR 1=1", Value: "paid"},
	}, DialectPostgres)
	if err == nil {
		t.Fatal("unsafe column should fail")
	}
}

func TestBuild_LeavesQueryUnchangedWithoutFilters(t *testing.T) {
	query, args, err := Build("SELECT id FROM orders", nil, DialectPostgres)
	if err != nil || query != "SELECT id FROM orders" || len(args) != 0 {
		t.Fatalf("Build() = (%q, %#v, %v)", query, args, err)
	}
}
