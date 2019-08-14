package noop

import (
	"github.com/Vinelab/tracing-go"
)

// Tracer is the tracing implementation for Zipkin. It should be initialized using NewTracer method.
type Tracer struct {
	extractionFormats map[string]tracing.Extractor
	injectionFormats  map[string]tracing.Injector
	rootSpan          tracing.Span
	currentSpan       tracing.Span
}

// NewTracer returns a new Zipkin tracer.
func NewTracer() *Tracer {
	return &Tracer{}
}

// StartSpan starts a new span based on a parent trace context. The context may come either from
// external source (extracted from HTTP request, AMQP message, etc., see Extract method)
// or received from another span in the service.
//
// If parent context does not contain a trace, a new trace will be implicitly created.
// Use EmptySpanContext to supply empty (nil) context.
func (tracer *Tracer) StartSpan(name string, spanCtx tracing.SpanContext) tracing.Span {
	var span *Span
	if tracer.rootSpan != nil {
		span = NewSpan(false)
	} else {
		span = NewSpan(true)
		tracer.rootSpan = span
	}

	tracer.currentSpan = span

	return span
}

// RootSpan retrieves the root span of the service
func (tracer *Tracer) RootSpan() tracing.Span {
	return tracer.rootSpan
}

// CurrentSpan retrieves the most recently activated span.
func (tracer *Tracer) CurrentSpan() tracing.Span {
	return tracer.currentSpan
}

// UUID retrieves unique identifier associated with a root span
func (tracer *Tracer) UUID() string {
	return ""
}

// EmptySpanContext return empty span context for creating spans
func (tracer *Tracer) EmptySpanContext() tracing.SpanContext {
	return NewSpanContext()
}

// Extract deserializes span context from from a given carrier using the format descriptor
// that tells tracer how to decode it from the carrier parameters
func (tracer *Tracer) Extract(carrier interface{}, format string) (tracing.SpanContext, error) {
	return NewSpanContext(), nil
}

// Inject implicitly serializes current span context using the format descriptor that
// tells how to encode trace info in the carrier parameters
func (tracer *Tracer) Inject(carrier interface{}, format string) error {
	span := tracer.currentSpan
	if span == nil {
		return tracing.ErrMissingTraceContext
	}

	return nil
}

// InjectContext serializes specified span context into a given carrier using the format descriptor
// that tells how to encode trace info in the carrier parameters
func (tracer *Tracer) InjectContext(carrier interface{}, format string, spanCtx tracing.SpanContext) error {
	return nil
}

// RegisterExtractionFormat register extractor implementation for given format string
func (tracer *Tracer) RegisterExtractionFormat(format string, extractor tracing.Extractor) {
	//
}

// RegisterInjectionFormat register injector implementation for given format string
func (tracer *Tracer) RegisterInjectionFormat(format string, injector tracing.Injector) {
	//
}

// Flush may flush any pending spans to the transport and reset the state of the tracer.
// Make sure this method is always called after the request is finished.
func (tracer *Tracer) Flush() {
	tracer.rootSpan = nil
	tracer.currentSpan = nil
}

// Close does a clean shutdown of the reporter, sending any traces that may be buffered in memory.
// This is especially useful for command-line tools that enable tracing,
// as well as for the long-running apps that support graceful shutdown.
//
// It goes without saying, but you cannot send anymore spans after calling Close,
// so you should only run this once during the lifecycle of the program.
func (tracer *Tracer) Close() error {
	return nil
}
