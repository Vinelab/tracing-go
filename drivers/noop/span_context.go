package noop

// SpanContext holds the context of a Span. It should be initialized using NewSpanContext method.
type SpanContext struct {
	//
}

// NewSpanContext returns a new SpanContext
func NewSpanContext() *SpanContext {
	return &SpanContext{}
}

// RawContext returns underlying (original) span context.
func (spanCtx *SpanContext) RawContext() interface{} {
	return nil
}
