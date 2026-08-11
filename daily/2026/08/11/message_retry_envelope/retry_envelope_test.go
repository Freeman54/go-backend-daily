package retryenvelope

import (
	"errors"
	"testing"
	"time"
)

func TestDecide(t *testing.T) {
	now := time.Date(2026, 8, 11, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		env  Envelope
		want Decision
		err  error
	}{
		{"retry", Envelope{1, now.Add(-time.Minute)}, Retry, nil},
		{"attempt exhausted", Envelope{3, now.Add(-time.Minute)}, DeadLetter, nil},
		{"age exhausted", Envelope{1, now.Add(-time.Hour)}, DeadLetter, nil},
		{"invalid attempt", Envelope{0, now}, DeadLetter, ErrInvalidAttempt},
		{"future timestamp", Envelope{1, now.Add(time.Minute)}, DeadLetter, ErrClockSkew},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decide(now, tt.env, 3, 30*time.Minute, 5*time.Second)
			if got != tt.want || !errors.Is(err, tt.err) {
				t.Fatalf("Decide() = %v, %v", got, err)
			}
		})
	}
}

func TestDecideRejectsInvalidPolicyAndMissingTime(t *testing.T) {
	now := time.Now()
	if _, err := Decide(now, Envelope{1, now}, 0, time.Minute, 0); !errors.Is(err, ErrInvalidAttempt) {
		t.Fatalf("error = %v", err)
	}
	if _, err := Decide(now, Envelope{1, time.Time{}}, 3, time.Minute, 0); !errors.Is(err, ErrClockSkew) {
		t.Fatalf("error = %v", err)
	}
}
