package httprangeplanner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidRange    = errors.New("invalid range")
	ErrRangeNotSatisfy = errors.New("range not satisfiable")
)

type Range struct {
	Start int64
	End   int64
}

func (r Range) Length() int64 { return r.End - r.Start + 1 }

// Plan 解析单段 HTTP bytes Range。多段范围应由调用方拒绝或交给专门的 multipart 实现。
func Plan(header string, size int64) (Range, error) {
	if size < 0 || !strings.HasPrefix(header, "bytes=") {
		return Range{}, ErrInvalidRange
	}
	spec := strings.TrimPrefix(header, "bytes=")
	if spec == "" || strings.Contains(spec, ",") {
		return Range{}, ErrInvalidRange
	}
	left, right, ok := strings.Cut(spec, "-")
	if !ok {
		return Range{}, ErrInvalidRange
	}
	if size == 0 {
		return Range{}, ErrRangeNotSatisfy
	}
	if left == "" {
		suffix, err := parseNonNegative(right)
		if err != nil || suffix == 0 {
			return Range{}, ErrInvalidRange
		}
		if suffix > size {
			suffix = size
		}
		return Range{Start: size - suffix, End: size - 1}, nil
	}
	start, err := parseNonNegative(left)
	if err != nil {
		return Range{}, ErrInvalidRange
	}
	if start >= size {
		return Range{}, ErrRangeNotSatisfy
	}
	end := size - 1
	if right != "" {
		end, err = parseNonNegative(right)
		if err != nil || end < start {
			return Range{}, ErrInvalidRange
		}
		if end >= size {
			end = size - 1
		}
	}
	return Range{Start: start, End: end}, nil
}

func ContentRange(r Range, size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.Start, r.End, size)
}

func parseNonNegative(s string) (int64, error) {
	if s == "" || strings.TrimSpace(s) != s || strings.HasPrefix(s, "+") {
		return 0, ErrInvalidRange
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < 0 {
		return 0, ErrInvalidRange
	}
	return n, nil
}
