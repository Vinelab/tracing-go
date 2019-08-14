package zipkin

import (
	"log"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/drivers/zipkin/propagation"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/streadway/amqp"
)

// AMQPInjector manages trace injection into AMQP carrier
type AMQPInjector struct {
	//
}

// NewAMQPInjector returns the instance of AMQPInjector
func NewAMQPInjector() *AMQPInjector {
	return &AMQPInjector{}
}

// Inject serialises given SpanContext into given amqp.Publishing object
func (extractor *AMQPInjector) Inject(spanCtx tracing.SpanContext, carrier interface{}) error {
	msg, ok := carrier.(*amqp.Publishing)
	if !ok {
		log.Fatalf("Expected *amqp.Publishing, got %T", carrier)
	}

	rawCtx := spanCtx.RawContext()

	zipkinCtx, ok := rawCtx.(model.SpanContext)
	if !ok {
		log.Fatalf("Expected %T, got %T", model.SpanContext{}, rawCtx)
	}

	inject := propagation.InjectAMQP(msg)
	return inject(zipkinCtx)
}
