package ttlcap

import (
	"errors"
	"testing"
	"time"
)

func TestPolicyApply(t *testing.T) {
	p := Policy{Min: time.Second, Max: time.Hour, NegativeMax: time.Minute}
	tests := []struct {
		name     string
		desired  time.Duration
		negative bool
		want     time.Duration
	}{
		{"normal", 5 * time.Minute, false, 5 * time.Minute},
		{"minimum", time.Millisecond, false, time.Second},
		{"maximum", 2 * time.Hour, false, time.Hour},
		{"default maximum", 0, false, time.Hour},
		{"negative maximum", 5 * time.Minute, true, time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Apply(tt.desired, tt.negative)
			if err != nil || got != tt.want {
				t.Fatalf("Apply() = %v, %v; want %v, nil", got, err, tt.want)
			}
		})
	}
}

func TestPolicyApplyRejectsInvalidPolicy(t *testing.T) {
	policies := []Policy{
		{},
		{Min: time.Minute, Max: time.Second, NegativeMax: time.Second},
		{Min: time.Second, Max: time.Hour, NegativeMax: time.Millisecond},
		{Min: time.Second, Max: time.Minute, NegativeMax: time.Hour},
	}
	for _, p := range policies {
		if _, err := p.Apply(time.Second, false); !errors.Is(err, ErrInvalidPolicy) {
			t.Fatalf("Apply() error = %v, want ErrInvalidPolicy", err)
		}
	}
}
