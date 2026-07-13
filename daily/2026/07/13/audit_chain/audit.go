package auditchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var ErrBrokenChain = errors.New("audit chain verification failed")

type Event struct {
	ID       string
	Actor    string
	Action   string
	Resource string
	At       time.Time
	Checksum string
}

type Trail struct {
	mu     sync.Mutex
	events []Event
}

func (t *Trail) Append(event Event) Event {
	t.mu.Lock()
	defer t.mu.Unlock()

	prev := ""
	if len(t.events) > 0 {
		prev = t.events[len(t.events)-1].Checksum
	}
	event.Checksum = checksum(prev, event)
	t.events = append(t.events, event)
	return event
}

func (t *Trail) Events() []Event {
	t.mu.Lock()
	defer t.mu.Unlock()

	cloned := make([]Event, len(t.events))
	copy(cloned, t.events)
	return cloned
}

func (t *Trail) Verify() error {
	return VerifyEvents(t.Events())
}

func VerifyEvents(events []Event) error {
	prev := ""
	for _, event := range events {
		if event.Checksum != checksum(prev, event) {
			return ErrBrokenChain
		}
		prev = event.Checksum
	}
	return nil
}

func checksum(prev string, event Event) string {
	sum := sha256.Sum256([]byte(
		prev + "|" + event.ID + "|" + event.Actor + "|" + event.Action + "|" + event.Resource + "|" + event.At.UTC().Format(time.RFC3339Nano),
	))
	return hex.EncodeToString(sum[:])
}
