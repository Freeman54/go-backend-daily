package slowstartlimiter

type Limiter struct {
	current int
	max     int
	success int
	step    int
}

func New(max, step int) *Limiter {
	if max < 1 {
		max = 1
	}
	if step < 1 {
		step = 1
	}
	return &Limiter{
		current: 1,
		max:     max,
		step:    step,
	}
}

func (l *Limiter) Limit() int {
	return l.current
}

func (l *Limiter) OnSuccess() {
	l.success++
	if l.success < l.step {
		return
	}
	l.success = 0
	if l.current < l.max {
		l.current++
	}
}

func (l *Limiter) OnFailure() {
	l.success = 0
	l.current = 1
}
