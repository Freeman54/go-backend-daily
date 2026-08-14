package cookieprefix

import (
	"net/http"
	"testing"
)

func TestValidateAcceptsValidPrefixedCookies(t *testing.T) {
	tests := []*http.Cookie{
		{Name: "__Secure-session", Secure: true},
		{Name: "__Host-session", Secure: true, Path: "/"},
		{Name: "regular", Path: "/app"},
	}
	for _, cookie := range tests {
		if err := Validate(cookie); err != nil {
			t.Fatalf("Validate(%q) error: %v", cookie.Name, err)
		}
	}
}

func TestValidateRejectsInvalidHostCookie(t *testing.T) {
	tests := []*http.Cookie{
		{Name: "__Host-session", Path: "/"},
		{Name: "__Host-session", Secure: true, Path: "/app"},
		{Name: "__Host-session", Secure: true, Path: "/", Domain: "example.com"},
		{Name: "__Secure-session"},
		nil,
	}
	for _, cookie := range tests {
		if err := Validate(cookie); err == nil {
			t.Fatalf("Validate(%#v) should fail", cookie)
		}
	}
}
