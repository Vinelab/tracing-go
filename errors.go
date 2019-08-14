package tracing

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingTraceContext is returned when trying to implicitly inject SpanContext
	// and there is no active span
	ErrMissingTraceContext = errors.New("cannot obtain context because there is no active span")
)

// UnregisteredFormatError is returned when you tried to inject/extract trace context
// using unregistered format or when there is a mismatch between format and extractor
type UnregisteredFormatError struct {
	format string
	err    string
}

// NewUnregisteredFormatError returns instance of UnregisteredFormatError
func NewUnregisteredFormatError(err string, format string) *UnregisteredFormatError {
	return &UnregisteredFormatError{err: err, format: format}
}

// Error returns the string representation of the error
func (e *UnregisteredFormatError) Error() string {
	return fmt.Sprintf("%s %s", e.err, e.format)
}
