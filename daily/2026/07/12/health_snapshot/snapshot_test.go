package healthsnapshot

import "testing"

func TestBuildReadySnapshot(t *testing.T) {
	snapshot := Build([]Report{
		{Dependency: "db", Critical: true, Status: StatusHealthy},
		{Dependency: "redis", Critical: false, Status: StatusHealthy},
	})

	if !snapshot.Ready || snapshot.Mode != "ready" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
}

func TestBuildDegradedSnapshot(t *testing.T) {
	snapshot := Build([]Report{
		{Dependency: "db", Critical: true, Status: StatusHealthy},
		{Dependency: "search", Critical: false, Status: StatusUnhealthy},
	})

	if !snapshot.Ready || snapshot.Mode != "degraded" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
	if snapshot.UnhealthyCount != 1 {
		t.Fatalf("unhealthy count = %d", snapshot.UnhealthyCount)
	}
}

func TestBuildUnreadySnapshotOnCriticalFailure(t *testing.T) {
	snapshot := Build([]Report{
		{Dependency: "db", Critical: true, Status: StatusUnhealthy},
		{Dependency: "cache", Critical: false, Status: StatusDegraded},
	})

	if snapshot.Ready || snapshot.Mode != "unready" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
	if len(snapshot.CriticalFailures) != 1 || snapshot.CriticalFailures[0] != "db" {
		t.Fatalf("critical failures = %#v", snapshot.CriticalFailures)
	}
}
