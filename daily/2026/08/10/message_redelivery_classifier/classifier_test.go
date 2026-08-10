package redelivery

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestClassifierDecide(t *testing.T) {
	c := Classifier{MaxAttempts: 5, MaxAge: 10 * time.Minute}
	transient := errors.New("broker unavailable")
	tests := []struct {
		name    string
		err     error
		attempt int
		age     time.Duration
		want    Decision
	}{
		{"success", nil, 1, time.Minute, Ack},
		{"transient", transient, 2, time.Minute, Retry},
		{"wrapped permanent", fmt.Errorf("consume: %w", &PermanentError{Err: errors.New("bad payload")}), 1, time.Minute, DeadLetter},
		{"attempt exhausted", transient, 5, time.Minute, DeadLetter},
		{"too old", transient, 2, 10 * time.Minute, DeadLetter},
		{"invalid attempt", transient, 0, time.Minute, DeadLetter},
		{"invalid age", transient, 1, -time.Second, DeadLetter},
		{"invalid config", transient, 1, time.Second, DeadLetter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classifier := c
			if tt.name == "invalid config" {
				classifier.MaxAttempts = 0
			}
			if got := classifier.Decide(tt.err, tt.attempt, tt.age); got != tt.want {
				t.Fatalf("Decide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermanentError(t *testing.T) {
	cause := errors.New("invalid schema")
	err := &PermanentError{Err: cause}
	if !errors.Is(err, cause) || err.Error() == "" {
		t.Fatal("PermanentError must expose its cause and a message")
	}
}
