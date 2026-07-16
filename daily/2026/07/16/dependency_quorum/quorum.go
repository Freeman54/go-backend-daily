package dependencyquorum

type Dependency struct {
	Name      string
	Weight    int
	Mandatory bool
	Healthy   bool
}

type Result struct {
	Healthy      bool
	Granted      int
	Required     int
	FailedChecks []string
}

func Evaluate(deps []Dependency, required int) Result {
	result := Result{Required: required}
	for _, dep := range deps {
		if !dep.Healthy {
			result.FailedChecks = append(result.FailedChecks, dep.Name)
			if dep.Mandatory {
				return result
			}
			continue
		}
		result.Granted += dep.Weight
	}
	result.Healthy = result.Granted >= required
	return result
}
