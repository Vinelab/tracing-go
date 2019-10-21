package zipkin

import (
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
)

// GooglePubSubExtractor manages trace extraction from Google PubSub carrier
type GooglePubSubExtractor struct {
	TracerSetter
}

// NewGooglePubSubExtractor returns the instance of AMQPExtractor
func NewGooglePubSubExtractor() *GooglePubSubExtractor {
	return &GooglePubSubExtractor{}
}

// Extract deserializes SpanContext from amqp.Delivery object
func (extractor *GooglePubSubExtractor) Extract(carrier interface{}) (tracing.SpanContext, error) {
	msg, ok := carrier.(*pubsub.Message)

	if !ok {
		log.Fatalf("Expected *pubsub.Message, got %T", carrier)
	}

	rawCtx := extractor.Tracing.Extract(propagation.ExtractGooglePubSub(msg))
	return NewSpanContext(rawCtx), nil
}