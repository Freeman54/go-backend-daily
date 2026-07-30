package sqlsavepointplan

import (
	"errors"
	"regexp"
)

var safeName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,62}$`)

type Plan struct{ name string }

func New(name string) (Plan, error) {
	if !safeName.MatchString(name) {
		return Plan{}, errors.New("savepoint name must be a safe SQL identifier")
	}
	return Plan{name: name}, nil
}

func (p Plan) Begin() string    { return `SAVEPOINT "` + p.name + `"` }
func (p Plan) Rollback() string { return `ROLLBACK TO SAVEPOINT "` + p.name + `"` }
func (p Plan) Release() string  { return `RELEASE SAVEPOINT "` + p.name + `"` }
