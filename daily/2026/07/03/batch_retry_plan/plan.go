package batchretryplan

type Outcome int

const (
	OutcomeSuccess Outcome = iota
	OutcomeRetryable
	OutcomePermanent
)

type ItemResult struct {
	ID      string
	Outcome Outcome
	Attempt int
}

type Plan struct {
	AckIDs   []string
	RetryIDs []string
	DLQIDs   []string
}

func Build(results []ItemResult, maxAttempts int) Plan {
	plan := Plan{
		AckIDs:   make([]string, 0, len(results)),
		RetryIDs: make([]string, 0, len(results)),
		DLQIDs:   make([]string, 0, len(results)),
	}

	for _, item := range results {
		switch item.Outcome {
		case OutcomeSuccess:
			plan.AckIDs = append(plan.AckIDs, item.ID)
		case OutcomePermanent:
			plan.DLQIDs = append(plan.DLQIDs, item.ID)
		case OutcomeRetryable:
			if item.Attempt+1 >= maxAttempts {
				plan.DLQIDs = append(plan.DLQIDs, item.ID)
				continue
			}
			plan.RetryIDs = append(plan.RetryIDs, item.ID)
		}
	}

	return plan
}
