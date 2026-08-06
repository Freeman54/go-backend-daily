// Package http_vary_cache_key builds deterministic HTTP cache keys.
package http_vary_cache_key

import (
	"errors"
	"net/http"
	"net/textproto"
	"sort"
	"strings"
)

// Build combines request identity with the selected Vary headers.
func Build(req *http.Request, vary []string) (string, error) {
	if req == nil || req.URL == nil {
		return "", errors.New("request and URL are required")
	}
	names := append([]string(nil), vary...)
	for i, name := range names {
		if !validHeaderName(name) {
			return "", errors.New("invalid vary header name")
		}
		names[i] = strings.ToLower(name)
	}
	sort.Strings(names)
	parts := []string{req.Method, strings.ToLower(req.Host), req.URL.EscapedPath(), req.URL.Query().Encode()}
	for _, name := range names {
		values := req.Header.Values(textproto.CanonicalMIMEHeaderKey(name))
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		parts = append(parts, name+":"+strings.Join(values, ","))
	}
	return strings.Join(parts, "|"), nil
}

func validHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("!#$%&'*+-.^_`|~", r)) {
			return false
		}
	}
	return true
}
