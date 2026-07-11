package logsamplingbudget

import (
	"testing"
	"time"
)

func TestSamplerEnforcesPerKeyBudget(t *testing.T) {
	base := time.Unix(0, 0)
	sampler := New(2, time.Minute)
	sampler.SetClock(func() time.Time { return base })

	if !sampler.Allow("db_timeout") {
		t.Fatal("first event should pass")
	}
	if !sampler.Allow("db_timeout") {
		t.Fatal("second event should pass")
	}
	if sampler.Allow("db_timeout") {
		t.Fatal("third event should be blocked")
	}
}

func TestSamplerResetsAfterWindow(t *testing.T) {
	base := time.Unix(0, 0)
	now := base
	sampler := New(1, time.Minute)
	sampler.SetClock(func() time.Time { return now })

	if !sampler.Allow("redis_error") {
		t.Fatal("first event should pass")
	}
	if sampler.Allow("redis_error") {
		t.Fatal("second event in same window should fail")
	}

	now = base.Add(2 * time.Minute)
	if !sampler.Allow("redis_error") {
		t.Fatal("event after window reset should pass")
	}
}

func TestSamplerSeparatesKeys(t *testing.T) {
	sampler := New(1, time.Minute)

	if !sampler.Allow("mq_lag") {
		t.Fatal("mq_lag should pass")
	}
	if !sampler.Allow("db_lag") {
		t.Fatal("db_lag should have independent budget")
	}
}
