package tracing

// Extractor interface can be used to provide the Tracer with custom
// implementations to deserialize data from a given carrier
type Extractor interface {
	// Extract deserializes span context from given carrier
	Extract(carrier interface{}) (SpanContext, error)
}
