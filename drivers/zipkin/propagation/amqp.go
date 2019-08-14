package propagation

import (
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/streadway/amqp"
)

// ExtractAMQP will extract a span.Context from the AMQP Message if found in B3 header format.
func ExtractAMQP(msg *amqp.Delivery) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		headers := make(map[string]string)

		for k, v := range msg.Headers {
			v, ok := v.(string)
			if !ok {
				continue
			}
			headers[k] = v
		}

		var (
			traceIDHeader      = headers[b3.TraceID]
			spanIDHeader       = headers[b3.SpanID]
			parentSpanIDHeader = headers[b3.ParentSpanID]
			sampledHeader      = headers[b3.Sampled]
			flagsHeader        = headers[b3.Flags]
		)

		return b3.ParseHeaders(
			traceIDHeader, spanIDHeader, parentSpanIDHeader, sampledHeader,
			flagsHeader,
		)
	}
}

// InjectAMQP will inject a span.Context into a AMQP message
func InjectAMQP(msg *amqp.Publishing) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return b3.ErrEmptyContext
		}

		if msg.Headers == nil {
			msg.Headers = amqp.Table{}
		}

		if sc.Debug {
			msg.Headers[b3.Flags] = "1"
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// so don't also send "X-B3-Sampled: 1".
			if *sc.Sampled {
				msg.Headers[b3.Sampled] = "1"
			} else {
				msg.Headers[b3.Sampled] = "0"
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			msg.Headers[b3.TraceID] = sc.TraceID.String()
			msg.Headers[b3.SpanID] = sc.ID.String()
			if sc.ParentID != nil {
				msg.Headers[b3.ParentSpanID] = sc.ParentID.String()
			}
		}

		return nil
	}
}
