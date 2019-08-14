package zipkin

import (
	"github.com/Vinelab/tracing-go"
	openzipkin "github.com/openzipkin/zipkin-go"
)

// Extractor is used to assert that the extractor implements SetTracing method
type Extractor interface {
	Extract(carrier interface{}) (tracing.SpanContext, error)
	SetTracing(tracer *openzipkin.Tracer)
}

// TracerSetter is supposed to be embedded in every Zipkin extractor
// to provide access to underlying Zipkin tracer instance
type TracerSetter struct {
	tracing.Extractor
	Tracing *openzipkin.Tracer
}

// SetTracing sets the instance of Zipkin tracer (from the underlying instrumnetation) on the embedding type
func (embedding *TracerSetter) SetTracing(tracer *openzipkin.Tracer) {
	embedding.Tracing = tracer
}
