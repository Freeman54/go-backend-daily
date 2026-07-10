package errorenvelope

import (
	"errors"
	"strings"
)

type Code string

const (
	CodeInternal     Code = "internal"
	CodeConflict     Code = "conflict"
	CodeUnauthorized Code = "unauthorized"
)

type CodedError struct {
	Code    Code
	Message string
	Err     error
}

func (e *CodedError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *CodedError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type Envelope struct {
	Code    Code
	Message string
	Causes  []string
}

func Wrap(code Code, message string, err error) error {
	return &CodedError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Build(err error) Envelope {
	if err == nil {
		return Envelope{}
	}

	env := Envelope{
		Code:    CodeInternal,
		Message: "internal server error",
		Causes:  flatten(err),
	}

	var coded *CodedError
	if errors.As(err, &coded) {
		env.Code = coded.Code
		env.Message = coded.Message
	}

	return env
}

func flatten(err error) []string {
	parts := make([]string, 0, 4)
	for err != nil {
		var coded *CodedError
		if errors.As(err, &coded) && coded == err {
			parts = append(parts, coded.Message)
			err = coded.Unwrap()
			continue
		}
		parts = append(parts, err.Error())
		err = errors.Unwrap(err)
	}
	return dedupe(parts)
}

func dedupe(parts []string) []string {
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		out = append(out, part)
	}
	return out
}
