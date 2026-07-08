package webhookreplaywindow

import (
	"errors"
	"testing"
	"time"
)

func TestVerifyAcceptsFreshSignedWebhook(t *testing.T) {
	now := time.Unix(1_000, 0)
	verifier := New("top-secret", 2*time.Minute)
	body := []byte(`{"event":"paid"}`)
	signature := verifier.Sign(now, "nonce-1", body)

	if err := verifier.Verify(now, "nonce-1", body, signature, now.Add(time.Second)); err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
}

func TestVerifyRejectsTimestampSkew(t *testing.T) {
	now := time.Unix(1_000, 0)
	verifier := New("top-secret", time.Minute)
	body := []byte(`{"event":"paid"}`)
	timestamp := now.Add(-2 * time.Minute)
	signature := verifier.Sign(timestamp, "nonce-2", body)

	err := verifier.Verify(timestamp, "nonce-2", body, signature, now)
	if !errors.Is(err, ErrTimestampSkew) {
		t.Fatalf("expected ErrTimestampSkew, got %v", err)
	}
}

func TestVerifyRejectsReplayNonce(t *testing.T) {
	now := time.Unix(1_000, 0)
	verifier := New("top-secret", 2*time.Minute)
	body := []byte(`{"event":"paid"}`)
	signature := verifier.Sign(now, "nonce-3", body)

	if err := verifier.Verify(now, "nonce-3", body, signature, now); err != nil {
		t.Fatalf("first Verify returned error: %v", err)
	}
	err := verifier.Verify(now, "nonce-3", body, signature, now.Add(time.Second))
	if !errors.Is(err, ErrReplayNonce) {
		t.Fatalf("expected ErrReplayNonce, got %v", err)
	}
}

func TestVerifyRejectsBadSignature(t *testing.T) {
	now := time.Unix(1_000, 0)
	verifier := New("top-secret", 2*time.Minute)
	body := []byte(`{"event":"paid"}`)

	err := verifier.Verify(now, "nonce-4", body, "deadbeef", now)
	if !errors.Is(err, ErrBadSignature) {
		t.Fatalf("expected ErrBadSignature, got %v", err)
	}
}
