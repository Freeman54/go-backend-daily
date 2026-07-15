package poisonmessagetracker

type Decision string

const (
	DecisionAck        Decision = "ack"
	DecisionRetry      Decision = "retry"
	DecisionQuarantine Decision = "quarantine"
)

type Tracker struct {
	maxAttempts int
	poisonCodes map[string]struct{}
	attempts    map[string]int
}

func New(maxAttempts int, poisonCodes []string) *Tracker {
	set := make(map[string]struct{}, len(poisonCodes))
	for _, code := range poisonCodes {
		set[code] = struct{}{}
	}
	return &Tracker{
		maxAttempts: maxAttempts,
		poisonCodes: set,
		attempts:    make(map[string]int),
	}
}

func (t *Tracker) Succeed(messageID string) Decision {
	delete(t.attempts, messageID)
	return DecisionAck
}

func (t *Tracker) Fail(messageID, code string) Decision {
	if _, ok := t.poisonCodes[code]; ok {
		delete(t.attempts, messageID)
		return DecisionQuarantine
	}

	t.attempts[messageID]++
	if t.attempts[messageID] >= t.maxAttempts {
		delete(t.attempts, messageID)
		return DecisionQuarantine
	}
	return DecisionRetry
}

func (t *Tracker) Attempts(messageID string) int {
	return t.attempts[messageID]
}
