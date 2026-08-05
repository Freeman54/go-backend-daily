package sql_isolation_selector

import "database/sql"

type Requirement struct {
	ReadOnly               bool
	AllowNonRepeatableRead bool
	PreventWriteSkew       bool
}

func Select(r Requirement) *sql.TxOptions {
	level := sql.LevelRepeatableRead
	if r.AllowNonRepeatableRead {
		level = sql.LevelReadCommitted
	}
	if r.PreventWriteSkew {
		level = sql.LevelSerializable
	}
	return &sql.TxOptions{Isolation: level, ReadOnly: r.ReadOnly}
}
