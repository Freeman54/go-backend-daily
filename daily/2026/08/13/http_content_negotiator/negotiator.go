// Package contentnegotiator selects a response media type from an HTTP Accept header.
package contentnegotiator

import (
	"mime"
	"sort"
	"strconv"
	"strings"
)

// Choose returns the best supported media type. Ties follow the server preference order.
func Choose(accept string, supported []string) (string, bool) {
	if strings.TrimSpace(accept) == "" {
		if len(supported) == 0 {
			return "", false
		}
		return supported[0], true
	}
	type offer struct {
		value string
		q     float64
		order int
	}
	var offers []offer
	for i, raw := range strings.Split(accept, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(raw))
		if err != nil {
			continue
		}
		q := 1.0
		if text, ok := params["q"]; ok {
			parsed, err := strconv.ParseFloat(text, 64)
			if err != nil || parsed < 0 || parsed > 1 {
				continue
			}
			q = parsed
		}
		if q > 0 {
			offers = append(offers, offer{strings.ToLower(mediaType), q, i})
		}
	}
	sort.SliceStable(offers, func(i, j int) bool { return offers[i].q > offers[j].q })
	for _, candidate := range offers {
		for _, value := range supported {
			if matches(candidate.value, strings.ToLower(value)) {
				return value, true
			}
		}
	}
	return "", false
}

func matches(pattern, value string) bool {
	if pattern == "*/*" {
		return true
	}
	p := strings.Split(pattern, "/")
	v := strings.Split(value, "/")
	return len(p) == 2 && len(v) == 2 && p[0] == v[0] && (p[1] == "*" || p[1] == v[1])
}
