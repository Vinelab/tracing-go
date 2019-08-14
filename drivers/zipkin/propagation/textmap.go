package propagation

import (
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

// ExtractTextMap will extract a span.Context from the string map if found in B3 header format.
func ExtractTextMap(dict map[string]string) propagation.Extractor {
	return func() (*model.SpanContext, error) {
		var (
			traceIDHeader      = dict[b3.TraceID]
			spanIDHeader       = dict[b3.SpanID]
			parentSpanIDHeader = dict[b3.ParentSpanID]
			sampledHeader      = dict[b3.Sampled]
			flagsHeader        = dict[b3.Flags]
		)

		return b3.ParseHeaders(
			traceIDHeader, spanIDHeader, parentSpanIDHeader, sampledHeader,
			flagsHeader,
		)
	}
}

// InjectHTTP will inject a span.Context into a string map
func InjectTextMap(dict map[string]string) propagation.Injector {
	return func(sc model.SpanContext) error {
		if (model.SpanContext{}) == sc {
			return b3.ErrEmptyContext
		}

		if sc.Debug {
			dict[b3.Flags] = "1"
		} else if sc.Sampled != nil {
			// Debug is encoded as X-B3-Flags: 1. Since Debug implies Sampled,
			// so don't also send "X-B3-Sampled: 1".
			if *sc.Sampled {
				dict[b3.Sampled] = "1"
			} else {
				dict[b3.Sampled] = "0"
			}
		}

		if !sc.TraceID.Empty() && sc.ID > 0 {
			dict[b3.TraceID] = sc.TraceID.String()
			dict[b3.SpanID] = sc.ID.String()
			if sc.ParentID != nil {
				dict[b3.ParentSpanID] = sc.ParentID.String()
			}
		}

		return nil
	}
}
