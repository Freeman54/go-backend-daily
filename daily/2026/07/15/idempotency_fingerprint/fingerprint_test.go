package idempotencyfingerprint

import (
	"net/url"
	"testing"
)

func TestSumIsStableAcrossMapAndQueryOrder(t *testing.T) {
	reqA := Request{
		Method: "post",
		Path:   "/payments",
		Query: url.Values{
			"expand": {"ledger", "events"},
			"env":    {"prod"},
		},
		Headers: map[string]string{
			"X-Tenant":     "acme",
			"Content-Type": "application/json",
		},
		Body: `{"amount":100,"currency":"CNY"}`,
	}
	reqB := Request{
		Method: "POST",
		Path:   "/payments",
		Query: url.Values{
			"env":    {"prod"},
			"expand": {"events", "ledger"},
		},
		Headers: map[string]string{
			"content-type": "application/json",
			"x-tenant":     "acme",
		},
		Body: `{"amount":100,"currency":"CNY"}`,
	}

	if gotA, gotB := Sum(reqA), Sum(reqB); gotA != gotB {
		t.Fatalf("fingerprint mismatch: %s != %s", gotA, gotB)
	}
}

func TestSumChangesWhenBodyChanges(t *testing.T) {
	base := Request{
		Method: "POST",
		Path:   "/payments",
		Body:   `{"amount":100}`,
	}

	changed := base
	changed.Body = `{"amount":200}`

	if Sum(base) == Sum(changed) {
		t.Fatal("fingerprint should change when body changes")
	}
}
