package deliveryretrybudget

type Decision int

const (
	Retry Decision = iota
	Quarantine
)

type Budget struct {
	MaxAttempts int
}

func (b Budget) Decide(attempts int, isTransient bool) Decision {
	if !isTransient {
		return Quarantine
	}
	if b.MaxAttempts <= 0 {
		return Quarantine
	}
	if attempts >= b.MaxAttempts {
		return Quarantine
	}
	return Retry
}
