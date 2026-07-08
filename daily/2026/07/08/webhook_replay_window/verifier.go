package webhookreplaywindow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrTimestampSkew = errors.New("timestamp outside allowed window")
	ErrReplayNonce   = errors.New("nonce already used")
	ErrBadSignature  = errors.New("signature mismatch")
)

type Verifier struct {
	mu      sync.Mutex
	secret  []byte
	maxSkew time.Duration
	seen    map[string]time.Time
}

func New(secret string, maxSkew time.Duration) *Verifier {
	if maxSkew <= 0 {
		maxSkew = 5 * time.Minute
	}

	return &Verifier{
		secret:  []byte(secret),
		maxSkew: maxSkew,
		seen:    make(map[string]time.Time),
	}
}

func (v *Verifier) Sign(timestamp time.Time, nonce string, body []byte) string {
	return v.sign(timestamp.Unix(), nonce, body)
}

func (v *Verifier) Verify(timestamp time.Time, nonce string, body []byte, signature string, now time.Time) error {
	if absDuration(now.Sub(timestamp)) > v.maxSkew {
		return ErrTimestampSkew
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.seen[nonce]; ok {
		return ErrReplayNonce
	}

	expected := v.sign(timestamp.Unix(), nonce, body)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return ErrBadSignature
	}

	v.seen[nonce] = now
	v.gc(now)
	return nil
}

func (v *Verifier) sign(unix int64, nonce string, body []byte) string {
	mac := hmac.New(sha256.New, v.secret)
	_, _ = mac.Write([]byte(fmt.Sprintf("%d.%s.", unix, nonce)))
	_, _ = mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func (v *Verifier) gc(now time.Time) {
	for nonce, seenAt := range v.seen {
		if now.Sub(seenAt) > v.maxSkew {
			delete(v.seen, nonce)
		}
	}
}

func absDuration(value time.Duration) time.Duration {
	if value < 0 {
		return -value
	}
	return value
}
