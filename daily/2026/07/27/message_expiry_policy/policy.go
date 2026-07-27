package message_expiry_policy

import (
	"errors"
	"time"
)

var (
	ErrExpired          = errors.New("message expired")
	ErrInvalidTimestamp = errors.New("current time is before message creation")
)

// Policy 计算消息在多级转发后还可使用的剩余生存时间。
type Policy struct {
	maxAge time.Duration
}

func New(maxAge time.Duration) (*Policy, error) {
	if maxAge <= 0 {
		return nil, errors.New("max age must be positive")
	}
	return &Policy{maxAge: maxAge}, nil
}

func (p *Policy) RemainingTTL(createdAt, now time.Time) (time.Duration, error) {
	if now.Before(createdAt) {
		return 0, ErrInvalidTimestamp
	}
	remaining := p.maxAge - now.Sub(createdAt)
	if remaining <= 0 {
		return 0, ErrExpired
	}
	return remaining, nil
}
