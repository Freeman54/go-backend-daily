package partialbatchack

type Decision string

const (
	DecisionAck   Decision = "ack"
	DecisionRetry Decision = "retry"
	DecisionDrop  Decision = "drop"
)

type ItemResult struct {
	ID       string
	Decision Decision
}

type Plan struct {
	AckIDs   []string
	RetryIDs []string
	DropIDs  []string
}

func BuildPlan(results []ItemResult) Plan {
	plan := Plan{}
	for _, result := range results {
		switch result.Decision {
		case DecisionAck:
			plan.AckIDs = append(plan.AckIDs, result.ID)
		case DecisionRetry:
			plan.RetryIDs = append(plan.RetryIDs, result.ID)
		case DecisionDrop:
			plan.DropIDs = append(plan.DropIDs, result.ID)
		}
	}
	return plan
}

func (p Plan) RetryRatio() float64 {
	total := len(p.AckIDs) + len(p.RetryIDs) + len(p.DropIDs)
	if total == 0 {
		return 0
	}
	return float64(len(p.RetryIDs)) / float64(total)
}
