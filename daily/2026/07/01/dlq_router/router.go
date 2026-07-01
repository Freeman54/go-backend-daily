package dlqrouter

import "errors"

var (
	ErrTransient = errors.New("transient failure")
	ErrPoison    = errors.New("poison message")
)

type Decision string

const (
	DecisionAck   Decision = "ack"
	DecisionRetry Decision = "retry"
	DecisionDLQ   Decision = "dlq"
)

type Message struct {
	Attempts int
}

type Classifier struct {
	MaxAttempts int
}

func (c Classifier) Decide(message Message, err error) Decision {
	if err == nil {
		return DecisionAck
	}
	if errors.Is(err, ErrPoison) {
		return DecisionDLQ
	}

	limit := c.MaxAttempts
	if limit <= 0 {
		limit = 3
	}
	if errors.Is(err, ErrTransient) && message.Attempts < limit {
		return DecisionRetry
	}
	return DecisionDLQ
}
