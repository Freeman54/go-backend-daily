package sloburnrate

import "errors"

var ErrInvalidObjective = errors.New("invalid slo objective")

type Window struct {
	Total  int
	Errors int
}

func (w Window) ErrorRate() float64 {
	if w.Total == 0 {
		return 0
	}
	return float64(w.Errors) / float64(w.Total)
}

func BurnRate(window Window, objective float64) (float64, error) {
	if objective <= 0 || objective >= 1 {
		return 0, ErrInvalidObjective
	}

	errorBudget := 1 - objective
	return window.ErrorRate() / errorBudget, nil
}

func ShouldAlert(short Window, long Window, objective float64, shortThreshold float64, longThreshold float64) (bool, error) {
	shortBurn, err := BurnRate(short, objective)
	if err != nil {
		return false, err
	}
	longBurn, err := BurnRate(long, objective)
	if err != nil {
		return false, err
	}
	return shortBurn >= shortThreshold && longBurn >= longThreshold, nil
}
