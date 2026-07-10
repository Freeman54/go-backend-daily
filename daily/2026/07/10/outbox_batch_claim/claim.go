package outboxbatchclaim

import (
	"sort"
	"time"
)

type Message struct {
	ID           string
	Attempts     int
	AvailableAt  time.Time
	ClaimedUntil time.Time
}

func Claim(now time.Time, lease time.Duration, limit int, messages []Message) []Message {
	if lease <= 0 {
		lease = time.Second
	}
	if limit <= 0 {
		return nil
	}

	candidates := make([]Message, 0, len(messages))
	for _, msg := range messages {
		if msg.AvailableAt.After(now) {
			continue
		}
		if msg.ClaimedUntil.After(now) {
			continue
		}
		candidates = append(candidates, msg)
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if !candidates[i].AvailableAt.Equal(candidates[j].AvailableAt) {
			return candidates[i].AvailableAt.Before(candidates[j].AvailableAt)
		}
		if candidates[i].Attempts != candidates[j].Attempts {
			return candidates[i].Attempts < candidates[j].Attempts
		}
		return candidates[i].ID < candidates[j].ID
	})

	if len(candidates) > limit {
		candidates = candidates[:limit]
	}
	for i := range candidates {
		candidates[i].ClaimedUntil = now.Add(lease)
		candidates[i].Attempts++
	}
	return candidates
}
