package zipkin

import (
	"fmt"
	"time"

	"github.com/Vinelab/tracing-go"
	"github.com/openzipkin/zipkin-go"
)

// Span encapsulates the state of logical operation it represents
type Span struct {
	rawSpan zipkin.Span
	isRoot  bool
}

// NewSpan returns a new Span
func NewSpan(rawSpan zipkin.Span, isRoot bool) *Span {
	return &Span{rawSpan: rawSpan, isRoot: isRoot}
}

// SetName sets (overrides) the string name for the logical operation this span represents.
func (span *Span) SetName(name string) {
	span.rawSpan.SetName(name)
}

// Tag give your span context for search, viewing and analysis. For example,
// a key "your_app.version" would let you lookup spans by version.
func (span *Span) Tag(key string, value string) {
	var sanitizedValue string
	if len(value) > MaxTagLen {
		sanitizedValue = fmt.Sprintf("Value exceeds the maximum allowed length of %d bytes", MaxTagLen)
	} else {
		sanitizedValue = value
	}

	span.rawSpan.Tag(key, sanitizedValue)
}

// Finish notifies that operation has finished. Span duration is derived by subtracting the start
// timestamp from this, and set when appropriate.
func (span *Span) Finish() {
	span.rawSpan.Finish()
}

// Annotate associates an event that explains latency with a timestamp.
func (span *Span) Annotate(message string) {
	span.rawSpan.Annotate(time.Now(), message)
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
	return NewSpanContext(span.rawSpan.Context())
}
