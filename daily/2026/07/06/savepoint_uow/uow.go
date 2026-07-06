package savepointuow

import "fmt"

type Scope struct {
	name  string
	outer bool
}

type Planner struct {
	depth   int
	counter int
}

func (p *Planner) Enter() ([]string, Scope) {
	if p.depth == 0 {
		p.depth++
		return []string{"BEGIN"}, Scope{outer: true}
	}

	p.counter++
	name := fmt.Sprintf("sp_%d", p.counter)
	p.depth++
	return []string{"SAVEPOINT " + name}, Scope{name: name}
}

func (p *Planner) Leave(scope Scope, success bool) []string {
	if p.depth == 0 {
		return nil
	}
	p.depth--

	if scope.outer {
		if success {
			return []string{"COMMIT"}
		}
		return []string{"ROLLBACK"}
	}

	if success {
		return []string{"RELEASE SAVEPOINT " + scope.name}
	}
	return []string{"ROLLBACK TO SAVEPOINT " + scope.name, "RELEASE SAVEPOINT " + scope.name}
}
