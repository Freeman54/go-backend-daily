// Package methodpolicy centralizes HTTP method validation and Allow headers.
package methodpolicy

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

var ErrMethodNotAllowed = errors.New("method not allowed")

// Policy stores a normalized immutable set of allowed HTTP methods.
type Policy struct {
	allowed map[string]struct{}
	header  string
}

// New builds a policy and rejects invalid or empty method sets.
func New(methods ...string) (Policy, error) {
	allowed := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		method = strings.ToUpper(strings.TrimSpace(method))
		if method == "" || !validToken(method) {
			return Policy{}, errors.New("invalid HTTP method")
		}
		allowed[method] = struct{}{}
	}
	if len(allowed) == 0 {
		return Policy{}, errors.New("empty HTTP method policy")
	}
	list := make([]string, 0, len(allowed))
	for method := range allowed {
		list = append(list, method)
	}
	sort.Strings(list)
	return Policy{allowed: allowed, header: strings.Join(list, ", ")}, nil
}

// Check validates a method and returns the canonical Allow header on rejection.
func (p Policy) Check(method string) (string, error) {
	if _, ok := p.allowed[strings.ToUpper(method)]; ok {
		return "", nil
	}
	return p.header, ErrMethodNotAllowed
}

// Middleware rejects disallowed methods before invoking next.
func (p Policy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, err := p.Check(r.Method)
		if err != nil {
			w.Header().Set("Allow", allow)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validToken(value string) bool {
	for _, r := range value {
		if !(r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}
