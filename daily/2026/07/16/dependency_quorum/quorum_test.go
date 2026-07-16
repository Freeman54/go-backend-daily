package dependencyquorum

import (
	"reflect"
	"testing"
)

func TestEvaluateFailsWhenMandatoryDependencyIsDown(t *testing.T) {
	result := Evaluate([]Dependency{
		{Name: "mysql", Weight: 3, Mandatory: true, Healthy: false},
		{Name: "redis", Weight: 1, Healthy: true},
	}, 3)

	if result.Healthy {
		t.Fatal("expected unhealthy result")
	}
	if !reflect.DeepEqual(result.FailedChecks, []string{"mysql"}) {
		t.Fatalf("unexpected failed checks: %#v", result.FailedChecks)
	}
}

func TestEvaluatePassesWhenWeightMeetsThreshold(t *testing.T) {
	result := Evaluate([]Dependency{
		{Name: "mysql", Weight: 3, Mandatory: true, Healthy: true},
		{Name: "profile-cache", Weight: 1, Healthy: false},
		{Name: "search", Weight: 2, Healthy: true},
	}, 5)

	if !result.Healthy {
		t.Fatalf("expected healthy result: %#v", result)
	}
}

func TestEvaluateReturnsUnhealthyWhenWeightIsInsufficient(t *testing.T) {
	result := Evaluate([]Dependency{
		{Name: "mysql", Weight: 3, Mandatory: true, Healthy: true},
		{Name: "search", Weight: 1, Healthy: false},
	}, 5)

	if result.Healthy {
		t.Fatalf("expected unhealthy result: %#v", result)
	}
	if result.Granted != 3 {
		t.Fatalf("unexpected granted weight: %d", result.Granted)
	}
}
