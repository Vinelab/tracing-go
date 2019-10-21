package zipkin

import (
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
	"github.com/openzipkin/zipkin-go/model"
)

// GooglePubSubInjector manages trace injection into Google Cloud PubSub carrier
type GooglePubSubInjector struct {
	//
}

// NewGooglePubSubInjector returns the instance of GooglePubSubInjector
func NewGooglePubSubInjector() *GooglePubSubInjector {
	return &GooglePubSubInjector{}
}

// Inject serialises given SpanContext into given amqp.Publishing object
func (extractor *GooglePubSubInjector) Inject(spanCtx tracing.SpanContext, carrier interface{}) error {
	msg, ok := carrier.(*pubsub.Message)
	if !ok {
		log.Fatalf("Expected *pubsub.Message, got %T", carrier)
	}

	rawCtx := spanCtx.RawContext()

	zipkinCtx, ok := rawCtx.(model.SpanContext)
	if !ok {
		log.Fatalf("Expected %T, got %T", model.SpanContext{}, rawCtx)
	}

	inject := propagation.InjectGooglePubSub(msg)
	return inject(zipkinCtx)
}
