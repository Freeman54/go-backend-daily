package tracingtraceparentvalidator

import (
	"encoding/hex"
	"errors"
	"strings"
)

var ErrInvalidTraceparent = errors.New("invalid traceparent")

type Parent struct {
	Version byte
	TraceID [16]byte
	SpanID  [8]byte
	Flags   byte
}

func (p Parent) Sampled() bool { return p.Flags&1 == 1 }

// Parse 校验 W3C traceparent 的基础格式。当前实现只接受 version 00。
func Parse(value string) (Parent, error) {
	var p Parent
	parts := strings.Split(value, "-")
	if len(parts) != 4 || len(parts[0]) != 2 || len(parts[1]) != 32 || len(parts[2]) != 16 || len(parts[3]) != 2 {
		return p, ErrInvalidTraceparent
	}
	version, err := decode(parts[0], 1)
	if err != nil || version[0] != 0 {
		return p, ErrInvalidTraceparent
	}
	traceID, err := decode(parts[1], 16)
	if err != nil || allZero(traceID) {
		return p, ErrInvalidTraceparent
	}
	spanID, err := decode(parts[2], 8)
	if err != nil || allZero(spanID) {
		return p, ErrInvalidTraceparent
	}
	flags, err := decode(parts[3], 1)
	if err != nil {
		return p, ErrInvalidTraceparent
	}
	p.Version, p.Flags = version[0], flags[0]
	copy(p.TraceID[:], traceID)
	copy(p.SpanID[:], spanID)
	return p, nil
}

func decode(s string, size int) ([]byte, error) {
	if s != strings.ToLower(s) {
		return nil, ErrInvalidTraceparent
	}
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != size {
		return nil, ErrInvalidTraceparent
	}
	return b, nil
}

func allZero(b []byte) bool {
	for _, value := range b {
		if value != 0 {
			return false
		}
	}
	return true
}
