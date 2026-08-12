// Package gapdetector detects missing offsets in an ordered message stream.
package gapdetector

import "sync"

// Detector tracks the next expected offset for one partition.
type Detector struct {
	mu   sync.Mutex
	next int64
}

func New(firstExpected int64) *Detector { return &Detector{next: firstExpected} }

// Observe classifies offset as duplicate, expected, or a gap without skipping missing data.
func (d *Detector) Observe(offset int64) Result {
	d.mu.Lock()
	defer d.mu.Unlock()
	result := Result{Expected: d.next, Offset: offset}
	switch {
	case offset < d.next:
		result.Kind = Duplicate
	case offset == d.next:
		result.Kind = Accepted
		d.next++
	default:
		result.Kind = Gap
		result.Missing = offset - d.next
	}
	return result
}

type Kind uint8

const (
	Accepted Kind = iota
	Duplicate
	Gap
)

type Result struct {
	Kind     Kind
	Expected int64
	Offset   int64
	Missing  int64
}
