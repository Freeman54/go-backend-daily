package sql_isolation_selector

import (
	"database/sql"
	"testing"
)

func TestSelect_MapsConsistencyRequirements(t *testing.T) {
	tests := []struct {
		name string
		in   Requirement
		want sql.IsolationLevel
	}{
		{"只读快照", Requirement{ReadOnly: true}, sql.LevelRepeatableRead},
		{"允许不可重复读", Requirement{AllowNonRepeatableRead: true}, sql.LevelReadCommitted},
		{"防止写偏差", Requirement{PreventWriteSkew: true}, sql.LevelSerializable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Select(tt.in); got.Isolation != tt.want || got.ReadOnly != tt.in.ReadOnly {
				t.Fatalf("got %+v", got)
			}
		})
	}
}
