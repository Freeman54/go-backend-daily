package canaryrollout

import "testing"

func TestEnabledHonorsLists(t *testing.T) {
	rule := Rule{
		Key:     "checkout-v2",
		Percent: 10,
		Allowlist: map[string]struct{}{
			"vip-user": {},
		},
		Denylist: map[string]struct{}{
			"blocked-user": {},
		},
	}

	if !Enabled(rule, "vip-user") {
		t.Fatal("allowlisted actor should be enabled")
	}
	if Enabled(rule, "blocked-user") {
		t.Fatal("denylisted actor should be disabled")
	}
}

func TestEnabledIsDeterministic(t *testing.T) {
	rule := Rule{Key: "search-rank", Percent: 25}

	first := Enabled(rule, "tenant-42")
	for i := 0; i < 10; i++ {
		if Enabled(rule, "tenant-42") != first {
			t.Fatal("expected deterministic decision for same actor")
		}
	}
}

func TestEnabledRespectsPercentEdges(t *testing.T) {
	if Enabled(Rule{Key: "x", Percent: 0}, "anyone") {
		t.Fatal("0 percent rollout should disable all actors")
	}
	if !Enabled(Rule{Key: "x", Percent: 100}, "anyone") {
		t.Fatal("100 percent rollout should enable all actors")
	}
}
