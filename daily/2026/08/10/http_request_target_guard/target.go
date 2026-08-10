// Package requesttargetguard validates and canonicalizes HTTP request targets.
package requesttargetguard

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

var (
	ErrTooLong = errors.New("request target exceeds budget")
	ErrInvalid = errors.New("invalid request target")
)

// Normalize accepts an HTTP origin-form request target and returns a canonical form.
func Normalize(target string, maxBytes int) (string, error) {
	if maxBytes <= 0 || len(target) > maxBytes {
		return "", ErrTooLong
	}
	if !strings.HasPrefix(target, "/") || strings.ContainsAny(target, "\\#") {
		return "", ErrInvalid
	}
	u, err := url.ParseRequestURI(target)
	if err != nil || u.IsAbs() || u.Host != "" {
		return "", ErrInvalid
	}
	decoded, err := url.PathUnescape(u.EscapedPath())
	if err != nil || strings.Contains(decoded, "\x00") || path.Clean(decoded) != decoded {
		return "", ErrInvalid
	}
	u.RawQuery = u.Query().Encode()
	return u.RequestURI(), nil
}
