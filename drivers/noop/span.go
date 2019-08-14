package noop

import (
	"github.com/Vinelab/tracing-go"
)

// Span encapsulates the state of logical operation it represents
type Span struct {
	isRoot bool
}

// NewSpan returns a new Span
func NewSpan(isRoot bool) *Span {
	return &Span{isRoot: isRoot}
}

// SetName sets (overrides) the string name for the logical operation this span represents.
func (span *Span) SetName(name string) {
	//
}

// Tag give your span context for search, viewing and analysis. For example,
// a key "your_app.version" would let you lookup spans by version.
func (span *Span) Tag(key string, value string) {
	//
}

// Finish notifies that operation has finished. Span duration is derived by subtracting the start
// timestamp from this, and set when appropriate.
func (span *Span) Finish() {
	//
}

// Annotate associates an event that explains latency with a timestamp.
func (span *Span) Annotate(message string) {
	//
}

// Log stores structured data. Despite this functionality being outlined in
// OpenTracing spec it's currently only supported in Jaeger
func (span *Span) Log(fields map[string]string) {
	//
}

// IsRoot tells whether the span is a root span
func (span *Span) IsRoot() bool {
	return span.isRoot
}

// Context retrieves SpanContext for this Span
func (span *Span) Context() tracing.SpanContext {
	return NewSpanContext()
}
