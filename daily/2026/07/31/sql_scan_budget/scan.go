package sqlscanbudget

import "errors"

var (
	ErrInvalidBudget  = errors.New("row and byte limits must be positive")
	ErrInvalidRowSize = errors.New("row size must not be negative")
	ErrBudgetExceeded = errors.New("scan budget exceeded")
)

// Budget 同时限制一次结果集扫描的行数和估算字节数。
type Budget struct {
	maxRows, maxBytes int
	rows, bytes       int
}

func New(maxRows, maxBytes int) (*Budget, error) {
	if maxRows <= 0 || maxBytes <= 0 {
		return nil, ErrInvalidBudget
	}
	return &Budget{maxRows: maxRows, maxBytes: maxBytes}, nil
}

func (b *Budget) Consume(rowBytes int) error {
	if rowBytes < 0 {
		return ErrInvalidRowSize
	}
	if b.rows+1 > b.maxRows || b.bytes+rowBytes > b.maxBytes {
		return ErrBudgetExceeded
	}
	b.rows++
	b.bytes += rowBytes
	return nil
}

func (b *Budget) Used() (rows, bytes int) {
	return b.rows, b.bytes
}
