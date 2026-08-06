// Package queue_backpressure_hysteresis implements a two-watermark controller.
package queue_backpressure_hysteresis

import (
	"errors"
	"sync"
)

type Controller struct {
	mu       sync.Mutex
	low      int
	high     int
	shedding bool
}

func New(low, high int) (*Controller, error) {
	if low < 0 || low >= high {
		return nil, errors.New("watermarks must satisfy 0 <= low < high")
	}
	return &Controller{low: low, high: high}, nil
}

// Observe returns whether producers should currently shed load.
func (c *Controller) Observe(depth int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.shedding && depth >= c.high {
		c.shedding = true
	} else if c.shedding && depth <= c.low {
		c.shedding = false
	}
	return c.shedding
}
