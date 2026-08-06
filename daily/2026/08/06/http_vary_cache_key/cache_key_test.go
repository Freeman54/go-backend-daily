package http_vary_cache_key

import (
	"net/http"
	"testing"
)

func TestBuild_CanonicalizesQueryAndVaryHeaders(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com/items?b=2&a=1", nil)
	req.Header.Add("Accept-Language", " zh-CN ")
	req.Header.Add("Accept-Language", "en-US")

	got, err := Build(req, []string{"accept-language"})
	if err != nil {
		t.Fatal(err)
	}
	want := "GET|api.example.com|/items|a=1&b=2|accept-language:zh-CN,en-US"
	if got != want {
		t.Fatalf("Build() = %q, want %q", got, want)
	}
}

func TestBuild_RejectsInvalidVaryName(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com", nil)
	if _, err := Build(req, []string{"X-Key\r\nInjected"}); err == nil {
		t.Fatal("Build() should reject an invalid header name")
	}
}
