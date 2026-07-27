package sql_value_clause

import "testing"

func TestBuild_CreatesPostgresPlaceholders(t *testing.T) {
	got, next, err := Build(2, 3, 4)
	if err != nil {
		t.Fatal(err)
	}
	want := "($4,$5),($6,$7),($8,$9)"
	if got != want || next != 10 {
		t.Fatalf("Build() = %q, %d; want %q, 10", got, next, want)
	}
}

func TestBuild_RejectsInvalidDimensions(t *testing.T) {
	for _, tc := range []struct{ columns, rows, start int }{
		{0, 1, 1}, {1, 0, 1}, {1, 1, 0},
	} {
		if _, _, err := Build(tc.columns, tc.rows, tc.start); err == nil {
			t.Fatalf("Build(%d, %d, %d) should fail", tc.columns, tc.rows, tc.start)
		}
	}
}
