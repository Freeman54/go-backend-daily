package slowquerysampler

import (
	"testing"
	"time"
)

func TestSamplerAlwaysLogsFailuresAndSlowQueries(t *testing.T) {
	sampler, err := New(150*time.Millisecond, 5)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if !sampler.ShouldLog("SELECT users", 20*time.Millisecond, true) {
		t.Fatal("expected failed query to be logged")
	}
	if !sampler.ShouldLog("SELECT users", 160*time.Millisecond, false) {
		t.Fatal("expected slow query to be logged")
	}
}

func TestSamplerUsesDeterministicSamplingForFastSuccessfulQueries(t *testing.T) {
	sampler, err := New(150*time.Millisecond, 1)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	first := sampler.ShouldLog("SELECT hot_table", 30*time.Millisecond, false)
	second := sampler.ShouldLog("SELECT hot_table", 30*time.Millisecond, false)
	if first != second {
		t.Fatal("sampling should be deterministic for the same operation")
	}

	samplerAll, err := New(150*time.Millisecond, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if !samplerAll.ShouldLog("SELECT hot_table", 30*time.Millisecond, false) {
		t.Fatal("100% sample rate should log all successful fast queries")
	}
}

func TestSamplerRejectsInvalidConfig(t *testing.T) {
	if _, err := New(0, 10); err == nil {
		t.Fatal("expected validation error for threshold")
	}
	if _, err := New(time.Second, 0); err == nil {
		t.Fatal("expected validation error for sampleRate")
	}
	if _, err := New(time.Second, 101); err == nil {
		t.Fatal("expected validation error for sampleRate")
	}
}
