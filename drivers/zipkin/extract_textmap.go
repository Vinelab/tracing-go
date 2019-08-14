package zipkin

import (
	"log"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
)

// TextMapExtractor manages trace extraction from TextMap carrier
type TextMapExtractor struct {
	TracerSetter
}

// NewTextMapExtractor returns the instance of TextMapExtractor
func NewTextMapExtractor() *TextMapExtractor {
	return &TextMapExtractor{}
}

// Extract deserializes SpanContext from a map
func (extractor *TextMapExtractor) Extract(carrier interface{}) (tracing.SpanContext, error) {
	textMap, ok := carrier.(map[string]string)

	if !ok {
		log.Fatalf("Expected map[string]string, got %T", carrier)
	}

	rawCtx := extractor.Tracing.Extract(propagation.ExtractTextMap(textMap))
	return NewSpanContext(rawCtx), nil
}
