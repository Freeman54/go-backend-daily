package outboxbatchclaim

import (
	"testing"
	"time"
)

func TestClaimSkipsUnavailableAndClaimedMessages(t *testing.T) {
	now := time.Unix(100, 0)
	messages := []Message{
		{ID: "a", AvailableAt: now.Add(-time.Second)},
		{ID: "b", AvailableAt: now.Add(time.Second)},
		{ID: "c", AvailableAt: now.Add(-time.Second), ClaimedUntil: now.Add(time.Second)},
	}

	got := Claim(now, 5*time.Second, 10, messages)
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected only message a, got %#v", got)
	}
	if got[0].Attempts != 1 {
		t.Fatalf("expected attempts increment, got %d", got[0].Attempts)
	}
}

func TestClaimOrdersByAvailabilityAttemptsAndID(t *testing.T) {
	now := time.Unix(100, 0)
	messages := []Message{
		{ID: "c", Attempts: 2, AvailableAt: now.Add(-time.Second)},
		{ID: "a", Attempts: 1, AvailableAt: now.Add(-time.Second)},
		{ID: "b", Attempts: 0, AvailableAt: now.Add(-2 * time.Second)},
	}

	got := Claim(now, 3*time.Second, 2, messages)
	if len(got) != 2 {
		t.Fatalf("expected 2 claimed messages, got %d", len(got))
	}
	if got[0].ID != "b" || got[1].ID != "a" {
		t.Fatalf("unexpected claim order: %#v", got)
	}
	if want := now.Add(3 * time.Second); !got[0].ClaimedUntil.Equal(want) {
		t.Fatalf("expected claimed until %v, got %v", want, got[0].ClaimedUntil)
	}
}
