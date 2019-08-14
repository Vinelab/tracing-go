package zipkin

// SpanContext holds the context of a Span. It should be initialized using NewSpanContext method.
type SpanContext struct {
	rawCtx interface{}
}

// NewSpanContext returns a new SpanContext
func NewSpanContext(rawCtx interface{}) *SpanContext {
	return &SpanContext{rawCtx: rawCtx}
}

// RawContext returns underlying (original) span context.
func (spanCtx *SpanContext) RawContext() interface{} {
	return spanCtx.rawCtx
}
