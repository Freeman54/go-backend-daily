package message_schema_gate

import "errors"

type Decision int

const (
	Accept Decision = iota
	Obsolete
	Future
)

type Gate struct{ min, max int }

func New(minVersion, maxVersion int) (Gate, error) {
	if minVersion <= 0 || maxVersion < minVersion {
		return Gate{}, errors.New("invalid schema version range")
	}
	return Gate{min: minVersion, max: maxVersion}, nil
}

func (g Gate) Check(version int) Decision {
	if version < g.min {
		return Obsolete
	}
	if version > g.max {
		return Future
	}
	return Accept
}
