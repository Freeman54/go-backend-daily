package exemplarsampler

import "testing"

func TestSampler_SelectsStableFractionByTraceID(t *testing.T) {
	s := New(25)
	first := s.Sample("trace-42")
	for i := 0; i < 10; i++ {
		if got := s.Sample("trace-42"); got != first {
			t.Fatalf("Sample() changed from %v to %v", first, got)
		}
	}
	if New(0).Sample("trace-42") {
		t.Fatal("0%% sampler should reject every trace")
	}
	if !New(100).Sample("trace-42") {
		t.Fatal("100%% sampler should accept every trace")
	}
}

func TestSampler_RejectsEmptyTraceAndClampsRate(t *testing.T) {
	if New(50).Sample("") {
		t.Fatal("empty trace ID should not produce an exemplar")
	}
	if New(-1).Rate() != 0 {
		t.Fatal("negative rate should clamp to 0")
	}
	if New(101).Rate() != 100 {
		t.Fatal("rate above 100 should clamp to 100")
	}
}
