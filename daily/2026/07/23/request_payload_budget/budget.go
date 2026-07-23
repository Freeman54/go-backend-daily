package requestpayloadbudget

import (
	"errors"
	"fmt"
	"io"
)

var ErrPayloadTooLarge = errors.New("request payload too large")

func Read(source io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("max bytes must be positive")
	}
	data, err := io.ReadAll(io.LimitReader(source, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read payload: %w", err)
	}
	if int64(len(data)) > maxBytes {
		return nil, ErrPayloadTooLarge
	}
	return data, nil
}
