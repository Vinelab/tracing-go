package propagation

import (
	"cloud.google.com/go/pubsub"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

// ExtractGooglePubSub will extract a span.Context from the Google Cloud PubSub message if found in B3 header format.
func ExtractGooglePubSub(msg *pubsub.Message) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		var (
			traceIDHeader      = msg.Attributes[b3.TraceID]
			spanIDHeader       = msg.Attributes[b3.SpanID]
			parentSpanIDHeader = msg.Attributes[b3.ParentSpanID]
			sampledHeader      = msg.Attributes[b3.Sampled]
			flagsHeader        = msg.Attributes[b3.Flags]
		)

		return b3.ParseHeaders(
			traceIDHeader, spanIDHeader, parentSpanIDHeader, sampledHeader,
			flagsHeader,
		)
	}
}

// InjectGooglePubSub will inject a span.Context into a Google Cloud PubSub message
func InjectGooglePubSub(msg *pubsub.Message) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return b3.ErrEmptyContext
		}

		if sc.Debug {
			msg.Attributes[b3.Flags] = "1"
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// so don't also send "X-B3-Sampled: 1".
			if *sc.Sampled {
				msg.Attributes[b3.Sampled] = "1"
			} else {
				msg.Attributes[b3.Sampled] = "0"
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			msg.Attributes[b3.TraceID] = sc.TraceID.String()
			msg.Attributes[b3.SpanID] = sc.ID.String()
			if sc.ParentID != nil {
				msg.Attributes[b3.ParentSpanID] = sc.ParentID.String()
			}
		}

		return nil
	}
}