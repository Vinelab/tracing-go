package tracing

// Span interface is returned by Tracer.StartSpan().
// You can use it to provide your own custom implementation
type Span interface {
	// SetName sets (overrides) the string name for the logical operation this span represents.
	SetName(name string)

	// Tag give your span context for search, viewing and analysis. For example,
	// a key "your_app.version" would let you lookup spans by version.
	Tag(key string, value string)

	// Finish notifies that operation has finished. Span duration is derived by subtracting the start
	// timestamp from this, and set when appropriate.
	Finish()

	// Annotate associates an event that explains latency with a timestamp.
	Annotate(message string)

	// Log stores structured data. Despite this functionality being outlined in
	// OpenTracing spec it's currently only supported in Jaeger
	Log(fields map[string]string)

	// IsRoot tells whether the span is a root span
	IsRoot() bool

	// Context retrieves SpanContext for this Span
	Context() SpanContext
}
