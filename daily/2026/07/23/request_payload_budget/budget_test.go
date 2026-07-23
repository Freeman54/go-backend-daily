package requestpayloadbudget

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestReadAcceptsPayloadWithinBudget(t *testing.T) {
	got, err := Read(bytes.NewBufferString("hello"), 5)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("Read() = %q, want hello", got)
	}
}

func TestReadRejectsPayloadOverBudget(t *testing.T) {
	_, err := Read(bytes.NewBufferString("hello!"), 5)
	if !errors.Is(err, ErrPayloadTooLarge) {
		t.Fatalf("Read() error = %v, want %v", err, ErrPayloadTooLarge)
	}
}

func TestReadRejectsInvalidBudget(t *testing.T) {
	if _, err := Read(bytes.NewBufferString(""), 0); err == nil {
		t.Fatal("Read() expected invalid budget error")
	}
}

func TestReadPropagatesSourceError(t *testing.T) {
	sentinel := errors.New("source failed")
	_, err := Read(errorReader{err: sentinel}, 10)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Read() error = %v, want %v", err, sentinel)
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}

var _ io.Reader = errorReader{}
