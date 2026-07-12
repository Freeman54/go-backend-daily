package healthsnapshot

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

type Report struct {
	Dependency string
	Critical   bool
	Status     Status
	Message    string
}

type Snapshot struct {
	Mode             string
	Ready            bool
	HealthyCount     int
	DegradedCount    int
	UnhealthyCount   int
	CriticalFailures []string
}

func Build(reports []Report) Snapshot {
	snapshot := Snapshot{
		Mode:  "ready",
		Ready: true,
	}

	for _, report := range reports {
		switch report.Status {
		case StatusHealthy:
			snapshot.HealthyCount++
		case StatusDegraded:
			snapshot.DegradedCount++
		case StatusUnhealthy:
			snapshot.UnhealthyCount++
			if report.Critical {
				snapshot.CriticalFailures = append(snapshot.CriticalFailures, report.Dependency)
			}
		}
	}

	if len(snapshot.CriticalFailures) > 0 {
		snapshot.Mode = "unready"
		snapshot.Ready = false
		return snapshot
	}

	if snapshot.DegradedCount > 0 || snapshot.UnhealthyCount > 0 {
		snapshot.Mode = "degraded"
	}

	return snapshot
}
