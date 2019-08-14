package tracing

// Tracer interface can be used to create your own tracing driver
// that wraps a lower level instrumentation
type Tracer interface {
	// StartSpan starts a new span based on a parent trace context. The context may come either from
	// external source (extracted from HTTP request, AMQP message, etc., see Extract method)
	// or received from another span in the service.
	//
	// If parent context does not contain a trace, a new trace will be implicitly created.
	// Use EmptySpanContext to supply empty (nil) context.
	StartSpan(name string, spanCtx SpanContext) Span

	// RootSpan retrieves the root span of the service
	RootSpan() Span

	// CurrentSpan retrieves the most recently activated span.
	CurrentSpan() Span

	// UUID retrieves unique identifier associated with a root span
	UUID() string

	// EmptySpanContext return empty span context for creating spans
	EmptySpanContext() SpanContext

	// Extract deserializes span context from from a given carrier using the format descriptor
	// that tells tracer how to decode it from the carrier parameters
	Extract(carrier interface{}, format string) (SpanContext, error)

	// Inject implicitly serializes current span context using the format descriptor that
	// tells how to encode trace info in the carrier parameters
	Inject(carrier interface{}, format string) error

	// InjectContext serializes specified span context into a given carrier using the format descriptor
	// that tells how to encode trace info in the carrier parameters
	InjectContext(carrier interface{}, format string, spanCtx SpanContext) error

	// RegisterExtractionFormat register extractor implementation for given format string
	RegisterExtractionFormat(format string, extractor Extractor)

	// RegisterInjectionFormat register injector implementation for given format string
	RegisterInjectionFormat(format string, injector Injector)

	// Flush may flush any pending spans to the transport and reset the state of the tracer.
	// Make sure this method is always called after the request is finished.
	Flush()

	// Close does a clean shutdown of the reporter, sending any traces that may be buffered in memory.
	// This is especially useful for command-line tools that enable tracing,
	// as well as for the long-running apps that support graceful shutdown.
	//
	// It goes without saying, but you cannot send anymore spans after calling Close,
	// so you should only run this once during the lifecycle of the program.
	Close() error
}
