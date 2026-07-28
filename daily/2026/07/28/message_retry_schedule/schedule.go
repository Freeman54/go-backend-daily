package message_retry_schedule

import (
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidPolicy  = errors.New("invalid retry policy")
	ErrInvalidAttempt = errors.New("attempt must be positive")
)

type Policy struct {
	Base   time.Duration
	Max    time.Duration
	Jitter float64
}

// Delay 计算第 attempt 次失败后的指数退避时间。
func (p Policy) Delay(attempt int, sample float64) (time.Duration, error) {
	if p.Base <= 0 || p.Max < p.Base || p.Jitter < 0 || p.Jitter > 1 || sample < 0 || sample > 1 {
		return 0, ErrInvalidPolicy
	}
	if attempt <= 0 {
		return 0, ErrInvalidAttempt
	}
	exponent := attempt - 1
	delay := p.Max
	if exponent < 63 {
		multiplier := math.Pow(2, float64(exponent))
		candidate := time.Duration(float64(p.Base) * multiplier)
		if candidate > 0 && candidate < p.Max {
			delay = candidate
		}
	}
	factor := 1 + (2*sample-1)*p.Jitter
	return time.Duration(float64(delay) * factor), nil
}
