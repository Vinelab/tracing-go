package zipkin

import (
	"log"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
	"github.com/streadway/amqp"
)

// AMQPExtractor manages trace extraction from AMQP carrier
type AMQPExtractor struct {
	TracerSetter
}

// NewAMQPExtractor returns the instance of AMQPExtractor
func NewAMQPExtractor() *AMQPExtractor {
	return &AMQPExtractor{}
}

// Extract deserializes SpanContext from amqp.Delivery object
func (extractor *AMQPExtractor) Extract(carrier interface{}) (tracing.SpanContext, error) {
	msg, ok := carrier.(*amqp.Delivery)

	if !ok {
		log.Fatalf("Expected *amqp.Delivery, got %T", carrier)
	}

	rawCtx := extractor.Tracing.Extract(propagation.ExtractAMQP(msg))
	return NewSpanContext(rawCtx), nil
}
