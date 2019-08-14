package zipkin

import (
	"log"
	"net/http"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
	"github.com/openzipkin/zipkin-go/model"
)

// HTTPInjector manages trace injection into HTTP carrier
type HTTPInjector struct {
	//
}

// NewHTTPInjector returns the instance of HTTPInjector
func NewHTTPInjector() *HTTPInjector {
	return &HTTPInjector{}
}

// Inject serialises given SpanContext into a given http.Request object
func (extractor *HTTPInjector) Inject(spanCtx tracing.SpanContext, carrier interface{}) error {
	req, ok := carrier.(*http.Request)
	if !ok {
		log.Fatalf("Expected *http.Request, got %T", carrier)
	}

	rawCtx := spanCtx.RawContext()

	zipkinCtx, ok := rawCtx.(model.SpanContext)
	if !ok {
		log.Fatalf("Expected %T, got %T", model.SpanContext{}, rawCtx)
	}

	inject := propagation.InjectHTTP(req)
	return inject(zipkinCtx)
}
