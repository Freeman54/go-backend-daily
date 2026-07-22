package hedgedrequestgate

import (
	"fmt"
	"time"
)

type Gate struct {
	delay     time.Duration
	maxHedges int
}

type Request struct {
	started   time.Time
	delay     time.Duration
	maxHedges int
	issued    int
}

func New(delay time.Duration, maxHedges int) (*Gate, error) {
	if delay <= 0 {
		return nil, fmt.Errorf("delay must be positive")
	}
	if maxHedges <= 0 {
		return nil, fmt.Errorf("max hedges must be positive")
	}
	return &Gate{delay: delay, maxHedges: maxHedges}, nil
}

func (g *Gate) Start(now time.Time) *Request {
	return &Request{started: now, delay: g.delay, maxHedges: g.maxHedges}
}

func (r *Request) ShouldHedge(now time.Time) bool {
	if r.issued >= r.maxHedges {
		return false
	}
	next := r.started.Add(time.Duration(r.issued+1) * r.delay)
	if now.Before(next) {
		return false
	}
	r.issued++
	return true
}
