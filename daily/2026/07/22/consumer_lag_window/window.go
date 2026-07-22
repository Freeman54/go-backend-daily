package consumerlagwindow

import "fmt"

type Snapshot struct {
	Count   int
	Max     int64
	Average int64
}

type Window struct {
	values []int64
	next   int
	full   bool
}

func New(size int) (*Window, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be positive")
	}
	return &Window{values: make([]int64, size)}, nil
}

func (w *Window) Add(lag int64) error {
	if lag < 0 {
		return fmt.Errorf("lag must not be negative")
	}
	w.values[w.next] = lag
	w.next++
	if w.next == len(w.values) {
		w.next = 0
		w.full = true
	}
	return nil
}

func (w *Window) Snapshot() Snapshot {
	count := w.next
	if w.full {
		count = len(w.values)
	}
	if count == 0 {
		return Snapshot{}
	}
	var sum, max int64
	for _, value := range w.values[:count] {
		sum += value
		if value > max {
			max = value
		}
	}
	return Snapshot{Count: count, Max: max, Average: sum / int64(count)}
}
