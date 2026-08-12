// Package responsebudget limits HTTP response bodies before they consume unbounded memory or bandwidth.
package responsebudget

import (
	"errors"
	"net/http"
)

var ErrResponseTooLarge = errors.New("response body exceeds budget")

// Writer wraps an http.ResponseWriter and rejects writes that exceed Limit.
// It is intended for buffered handlers that can turn the error into a controlled response.
type Writer struct {
	http.ResponseWriter
	Limit   int64
	written int64
}

// Write accepts p only when the complete write remains within the configured budget.
func (w *Writer) Write(p []byte) (int, error) {
	if w.Limit < 0 || int64(len(p)) > w.Limit-w.written {
		return 0, ErrResponseTooLarge
	}
	n, err := w.ResponseWriter.Write(p)
	w.written += int64(n)
	return n, err
}

// Written reports the bytes successfully delegated to the underlying writer.
func (w *Writer) Written() int64 { return w.written }
