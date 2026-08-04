package sql_limit_policy

import "testing"

func TestPolicy_NormalizesRequestedLimit(t *testing.T) {
	p, err := New(20, 100)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct{ requested, want int }{{0, 20}, {-1, 20}, {1, 1}, {100, 100}, {101, 100}}
	for _, tt := range tests {
		if got := p.Normalize(tt.requested); got != tt.want {
			t.Errorf("Normalize(%d) = %d, want %d", tt.requested, got, tt.want)
		}
	}
}

func TestNew_RejectsInvalidConfiguration(t *testing.T) {
	for _, tc := range [][2]int{{0, 100}, {20, 0}, {101, 100}} {
		if _, err := New(tc[0], tc[1]); err == nil {
			t.Errorf("New(%d, %d) error = nil", tc[0], tc[1])
		}
	}
}
