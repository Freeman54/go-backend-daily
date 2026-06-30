package readinessgate

import "sort"

type Dependency struct {
	Name  string
	Ready bool
}

type Snapshot struct {
	Ready   bool
	Reasons []string
}

type Gate struct {
	draining bool
	deps     map[string]bool
}

func New(names ...string) *Gate {
	deps := make(map[string]bool, len(names))
	for _, name := range names {
		deps[name] = false
	}
	return &Gate{deps: deps}
}

func (g *Gate) SetDependency(name string, ready bool) {
	if g.deps == nil {
		g.deps = make(map[string]bool)
	}
	g.deps[name] = ready
}

func (g *Gate) SetDraining(draining bool) {
	g.draining = draining
}

func (g *Gate) Snapshot() Snapshot {
	reasons := make([]string, 0)
	if g.draining {
		reasons = append(reasons, "instance is draining")
	}
	for name, ready := range g.deps {
		if !ready {
			reasons = append(reasons, name+" is not ready")
		}
	}
	sort.Strings(reasons)
	return Snapshot{
		Ready:   len(reasons) == 0,
		Reasons: reasons,
	}
}
