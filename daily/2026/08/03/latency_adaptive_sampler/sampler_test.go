package latency_adaptive_sampler

import (
	"testing"
	"time"
)

func TestSampler_AlwaysKeepsSlowAndErrorRequests(t *testing.T) {
	s := New(100*time.Millisecond, 10)
	if !s.Keep(101*time.Millisecond, false, 99) {
		t.Fatal("slow request should be kept")
	}
	if !s.Keep(time.Millisecond, true, 99) {
		t.Fatal("failed request should be kept")
	}
}

func TestSampler_SamplesFastSuccessDeterministically(t *testing.T) {
	s := New(100*time.Millisecond, 10)
	if !s.Keep(time.Millisecond, false, 20) || s.Keep(time.Millisecond, false, 21) {
		t.Fatal("one in ten fast successes should be kept by stable hash")
	}
}

func TestNew_ZeroFastRateKeepsFastSuccess(t *testing.T) {
	s := New(100*time.Millisecond, 0)
	if !s.Keep(time.Millisecond, false, 123) {
		t.Fatal("zero rate should be normalized to keep all requests")
	}
}
