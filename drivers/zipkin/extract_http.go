package zipkin

import (
	"log"
	"net/http"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
)

// HTTPExtractor manages trace extraction from HTTP carrier
type HTTPExtractor struct {
	TracerSetter
}

// NewHTTPExtractor returns the instance of HTTPExtractor
func NewHTTPExtractor() *HTTPExtractor {
	return &HTTPExtractor{}
}

// Extract deserializes SpanContext from http.Request object
func (extractor *HTTPExtractor) Extract(carrier interface{}) (tracing.SpanContext, error) {
	request, ok := carrier.(*http.Request)

	if !ok {
		log.Fatalf("Expected *http.Request, got %T", carrier)
	}

	rawCtx := extractor.Tracing.Extract(propagation.ExtractHTTP(request))
	return NewSpanContext(rawCtx), nil
}
