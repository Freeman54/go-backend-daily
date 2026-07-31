package problemdetailsanitizer

import "testing"

func TestSanitize_ExposesAllowedFieldsAndRedactsSecrets(t *testing.T) {
	input := Problem{
		Type:     "https://example.com/problems/invalid-order",
		Title:    "Invalid order",
		Status:   422,
		Detail:   "token=secret database host db.internal rejected card",
		Instance: "/orders/42",
	}
	got := Sanitize(input, []string{"rejected card"})
	if got.Detail != "rejected card" {
		t.Fatalf("Detail = %q, want safe public detail", got.Detail)
	}
	if got.Instance != "" {
		t.Fatalf("Instance = %q, want redacted", got.Instance)
	}
	if got.Type != input.Type || got.Title != input.Title || got.Status != input.Status {
		t.Fatalf("stable fields changed: %#v", got)
	}
}

func TestSanitize_FallsBackForUnknownDetailAndInvalidStatus(t *testing.T) {
	got := Sanitize(Problem{Status: 99, Detail: "internal stack"}, nil)
	if got.Status != 500 || got.Detail != "request could not be completed" {
		t.Fatalf("Sanitize() = %#v", got)
	}
}
