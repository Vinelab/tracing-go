package tracing

type SpanContext interface {
	// RawContext returns underlying (original) span context.
	RawContext() interface{}
}
