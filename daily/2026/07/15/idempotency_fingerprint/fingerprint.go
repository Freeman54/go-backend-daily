package idempotencyfingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers map[string]string
	Body    string
}

func CanonicalString(r Request) string {
	var b strings.Builder
	b.WriteString(strings.ToUpper(r.Method))
	b.WriteByte('\n')
	b.WriteString(r.Path)
	b.WriteByte('\n')
	b.WriteString(canonicalValues(r.Query))
	b.WriteByte('\n')
	b.WriteString(canonicalHeaders(r.Headers))
	b.WriteByte('\n')
	b.WriteString(r.Body)
	return b.String()
}

func Sum(r Request) string {
	sum := sha256.Sum256([]byte(CanonicalString(r)))
	return hex.EncodeToString(sum[:])
}

func canonicalValues(values url.Values) string {
	if len(values) == 0 {
		return ""
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		items := append([]string(nil), values[key]...)
		sort.Strings(items)
		parts = append(parts, key+"="+strings.Join(items, ","))
	}
	return strings.Join(parts, "&")
}

func canonicalHeaders(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}

	normalized := make(map[string]string, len(headers))
	for key, value := range headers {
		normalized[strings.ToLower(key)] = strings.TrimSpace(value)
	}

	keys := make([]string, 0, len(normalized))
	for key := range normalized {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+normalized[key])
	}
	return strings.Join(parts, "&")
}
