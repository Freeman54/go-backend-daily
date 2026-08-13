package contentnegotiator

import "testing"

func TestChoose(t *testing.T) {
	tests := []struct {
		name, accept, want string
		ok                 bool
	}{
		{"quality", "application/json;q=0.5, application/problem+json", "application/problem+json", true},
		{"wildcard", "text/*", "text/plain", true}, {"default", "", "application/json", true},
		{"unacceptable", "application/json;q=0", "", false}, {"invalid", "broken;=x", "", false},
	}
	supported := []string{"application/json", "application/problem+json", "text/plain"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Choose(tt.accept, supported)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("Choose() = %q, %v", got, ok)
			}
		})
	}
}

func TestChooseEmptySupported(t *testing.T) {
	if _, ok := Choose("", nil); ok {
		t.Fatal("expected no match")
	}
}
