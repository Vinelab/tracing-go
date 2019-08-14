package tracing

// Injector interface can be used to provide the Tracer with custom
// implementations to serialize data into a given carrier
type Injector interface {
	// Inject serializes span context into given carrier
	Inject(spanCtx SpanContext, carrier interface{}) error
}
