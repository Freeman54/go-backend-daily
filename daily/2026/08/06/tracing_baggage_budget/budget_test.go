package tracing_baggage_budget

import "testing"

func TestApply_PrioritizesAllowlistWithinBudget(t *testing.T) {
	got, dropped := Apply(map[string]string{
		"tenant": "acme",
		"region": "cn",
		"debug":  "verbose",
	}, []string{"tenant", "region", "debug"}, 2, 20)
	if len(got) != 2 || got["tenant"] != "acme" || got["region"] != "cn" {
		t.Fatalf("Apply() got %#v", got)
	}
	if dropped != 1 {
		t.Fatalf("dropped = %d, want 1", dropped)
	}
}

func TestApply_RejectsEntriesOverByteBudget(t *testing.T) {
	got, dropped := Apply(map[string]string{"trace": "123456"}, []string{"trace"}, 1, 8)
	if len(got) != 0 || dropped != 1 {
		t.Fatalf("Apply() = %#v, %d; want empty, 1", got, dropped)
	}
}
