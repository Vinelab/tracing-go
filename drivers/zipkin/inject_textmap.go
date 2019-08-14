package zipkin

import (
	"log"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
	"github.com/openzipkin/zipkin-go/model"
)

// TextMapInjector manages trace injection into TextMap carrier
type TextMapInjector struct {
	//
}

// NewTextMapInjector returns the instance of TextMapInjector
func NewTextMapInjector() *TextMapInjector {
	return &TextMapInjector{}
}

// Inject serialises given SpanContext into a given map
func (extractor *TextMapInjector) Inject(spanCtx tracing.SpanContext, carrier interface{}) error {
	textMap, ok := carrier.(*map[string]string)
	if !ok {
		log.Fatalf("Expected *map[string]string, got %T", carrier)
	}

	rawCtx := spanCtx.RawContext()

	zipkinCtx, ok := rawCtx.(model.SpanContext)
	if !ok {
		log.Fatalf("Expected %T, got %T", model.SpanContext{}, rawCtx)
	}

	inject := propagation.InjectTextMap(*textMap)
	return inject(zipkinCtx)
}
