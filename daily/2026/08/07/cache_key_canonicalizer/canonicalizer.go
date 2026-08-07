// Package cachekeycanonicalizer creates stable cache keys from named dimensions.
package cachekeycanonicalizer

import (
	"errors"
	"net/url"
	"strings"
)

// Build returns namespace followed by a deterministically ordered query string.
func Build(namespace string, parts map[string]string) (string, error) {
	if strings.TrimSpace(namespace) == "" {
		return "", errors.New("cache namespace is required")
	}
	values := make(url.Values, len(parts))
	for name, value := range parts {
		if name == "" {
			return "", errors.New("cache key part name is required")
		}
		values.Set(name, value)
	}
	encoded := values.Encode()
	if encoded == "" {
		return namespace, nil
	}
	return namespace + "?" + encoded, nil
}
