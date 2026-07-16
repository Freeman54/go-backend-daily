package queryshapehash

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	stringLiteralRE = regexp.MustCompile(`'([^']|'')*'`)
	numberLiteralRE = regexp.MustCompile(`\b\d+(\.\d+)?\b`)
	spaceRE         = regexp.MustCompile(`\s+`)
)

func Normalize(query string) string {
	normalized := strings.TrimSpace(query)
	normalized = strings.ToLower(normalized)
	normalized = stringLiteralRE.ReplaceAllString(normalized, "?")
	normalized = numberLiteralRE.ReplaceAllString(normalized, "?")
	normalized = spaceRE.ReplaceAllString(normalized, " ")
	return normalized
}

func Hash(query string) string {
	sum := sha1.Sum([]byte(Normalize(query)))
	return hex.EncodeToString(sum[:8])
}
