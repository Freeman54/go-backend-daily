package sloburnrate

import (
	"errors"
	"math"
	"testing"
)

func TestBurnRate(t *testing.T) {
	t.Parallel()

	got, err := BurnRate(Window{Total: 1000, Errors: 20}, 0.99)
	if err != nil {
		t.Fatalf("BurnRate() error = %v", err)
	}
	if math.Abs(got-2) > 1e-9 {
		t.Fatalf("BurnRate() = %v, want 2", got)
	}
}

func TestShouldAlertWhenShortAndLongWindowsBothHigh(t *testing.T) {
	t.Parallel()

	notify, err := ShouldAlert(
		Window{Total: 100, Errors: 20},
		Window{Total: 1000, Errors: 80},
		0.99,
		10,
		5,
	)
	if err != nil {
		t.Fatalf("ShouldAlert() error = %v", err)
	}
	if !notify {
		t.Fatalf("ShouldAlert() = false, want true")
	}
}

func TestShouldNotAlertWhenOnlyOneWindowIsHigh(t *testing.T) {
	t.Parallel()

	notify, err := ShouldAlert(
		Window{Total: 100, Errors: 20},
		Window{Total: 1000, Errors: 10},
		0.99,
		10,
		5,
	)
	if err != nil {
		t.Fatalf("ShouldAlert() error = %v", err)
	}
	if notify {
		t.Fatalf("ShouldAlert() = true, want false")
	}
}

func TestBurnRateRejectsInvalidObjective(t *testing.T) {
	t.Parallel()

	_, err := BurnRate(Window{Total: 10, Errors: 1}, 1)
	if !errors.Is(err, ErrInvalidObjective) {
		t.Fatalf("BurnRate() error = %v, want ErrInvalidObjective", err)
	}
}
