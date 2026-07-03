package tracesampler

import (
	"testing"
	"time"
)

func TestTraceSamplerAlwaysKeepsForcedErrorAndSlow(t *testing.T) {
	sampler := New(0, 200*time.Millisecond)

	for _, tc := range []struct {
		name string
		meta Meta
		want string
	}{
		{name: "forced", meta: Meta{TraceID: "a", Forced: true}, want: "forced"},
		{name: "error", meta: Meta{TraceID: "b", HasError: true}, want: "error"},
		{name: "slow", meta: Meta{TraceID: "c", Duration: 250 * time.Millisecond}, want: "slow"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			decision := sampler.Decide(tc.meta)
			if !decision.Sample || decision.Reason != tc.want {
				t.Fatalf("unexpected decision: %#v", decision)
			}
		})
	}
}

func TestTraceSamplerIsDeterministicForBaseSampling(t *testing.T) {
	sampler := New(300, time.Second)
	meta := Meta{TraceID: "trace-42", Duration: 10 * time.Millisecond}

	first := sampler.Decide(meta)
	second := sampler.Decide(meta)
	if first != second {
		t.Fatalf("expected deterministic sampling, got first=%#v second=%#v", first, second)
	}
}

func TestTraceSamplerDropsWhenBaseRateIsZero(t *testing.T) {
	sampler := New(0, time.Second)
	decision := sampler.Decide(Meta{TraceID: "trace-99", Duration: 10 * time.Millisecond})
	if decision.Sample || decision.Reason != "base-drop" {
		t.Fatalf("unexpected decision: %#v", decision)
	}
}
