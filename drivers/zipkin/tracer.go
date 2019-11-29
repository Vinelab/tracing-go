package zipkin

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/formats"
	"github.com/google/uuid"
	openzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

const (
	// MaxTagLen controls the maximum size of tag value in bytes
	MaxTagLen = 1048576
	// DefaultRequestTimeout sets maximum timeout for http request to send spans
	DefaultRequestTimeout = time.Second * 5
)

var (
	// ErrCollectorIPNotFound is returned if you used hostname in place of collector IP
	// and we weren't able to resolve the valid address from it
	ErrCollectorIPNotFound = errors.New("unable to resolve collector's IP address")
)

// Tracer is the tracing implementation for Zipkin. It should be initialized using NewTracer method.
type Tracer struct {
	tracing           *openzipkin.Tracer
	reporter          reporter.Reporter
	extractionFormats map[string]tracing.Extractor
	injectionFormats  map[string]tracing.Injector
	rootSpan          tracing.Span
	currentSpan       tracing.Span
	uuid              string
}

// TracerOptions is a configuration container to setup the Tracer.
type TracerOptions struct {
	// ServiceName is the name of application you're tracing
	// Required
	ServiceName string
	// Host
	// Required
	Host string
	// Port
	// Required
	Port string
	// UsesTraceID128Bit tells whether to use 128 bit trace IDs (32 characters in length as opposed to 16)
	// Defaults to false
	UsesTraceID128Bit bool
	// Reporter option allows to inject your own reporter for tests
	// Defaults to http reporter
	Reporter reporter.Reporter
	// Timeout sets maximum timeout for http request to send spans
	// Setting this to a too high value is not recommended because it
	// 	may degrade your system performance when collector is down.
	// Note that reporter will re-try after the first failure.
	// 	See this issue for more details: https://github.com/openzipkin/zipkin-go/issues/147
	RequestTimeout time.Duration
}

// NewTracer returns a new Zipkin tracer.
func NewTracer(opt TracerOptions) (*Tracer, error) {
	ipAddr, err := resolveCollectorIP(opt.Host)
	if err != nil {
		log.Printf("Unable to resolve collector's IP address from hostname %s: %s", opt.Host, err.Error())
	}
	opt.Host = ipAddr

	var timeout time.Duration
	if opt.RequestTimeout != 0 {
		timeout = opt.RequestTimeout
	} else {
		timeout = DefaultRequestTimeout
	}

	var rep reporter.Reporter
	if opt.Reporter != nil {
		rep = opt.Reporter
	} else {
		url := fmt.Sprintf("http://%s:%s/api/v2/spans", opt.Host, opt.Port)
		rep = httpreporter.NewReporter(url, httpreporter.Timeout(timeout))
	}

	endpoint, err := openzipkin.NewEndpoint(opt.ServiceName, fmt.Sprintf("%s:%s", opt.Host, opt.Port))
	if err != nil {
		return nil, err
	}

	trace, err := openzipkin.NewTracer(
		rep,
		openzipkin.WithLocalEndpoint(endpoint),
		openzipkin.WithTraceID128Bit(opt.UsesTraceID128Bit),
	)
	if err != nil {
		return nil, err
	}

	return &Tracer{
		tracing:           trace,
		reporter:          rep,
		extractionFormats: registerDefaultExtractionFormats(),
		injectionFormats:  registerDefaultInjectionFormats(),
	}, nil
}

// StartSpan starts a new span based on a parent trace context. The context may come either from
// external source (extracted from HTTP request, AMQP message, etc., see Extract method)
// or received from another span in the service.
//
// If parent context does not contain a trace, a new trace will be implicitly created.
// Use EmptySpanContext to supply empty (nil) context.
func (tracer *Tracer) StartSpan(name string, spanCtx tracing.SpanContext) tracing.Span {
	rawCtx := spanCtx.RawContext()
	parent, ok := rawCtx.(model.SpanContext)

	var rawSpan openzipkin.Span
	if ok {
		rawSpan = tracer.tracing.StartSpan(name, openzipkin.Parent(parent))
	} else {
		rawSpan = tracer.tracing.StartSpan(name)
	}

	var span *Span
	if tracer.rootSpan != nil {
		span = NewSpan(rawSpan, false)
	} else {
		span = NewSpan(rawSpan, true)
		tracer.rootSpan = span

		value, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		tracer.uuid = value.String()
		span.Tag("uuid", tracer.uuid)
	}

	tracer.currentSpan = span
	span.SetName(name)

	return span
}

// RootSpan retrieves the root span of the service
func (tracer *Tracer) RootSpan() tracing.Span {
	return tracer.rootSpan
}

// CurrentSpan retrieves the most recently activated span.
func (tracer *Tracer) CurrentSpan() tracing.Span {
	return tracer.currentSpan
}

// UUID retrieves unique identifier associated with a root span
func (tracer *Tracer) UUID() string {
	return tracer.uuid
}

// EmptySpanContext return empty span context for creating spans
func (tracer *Tracer) EmptySpanContext() tracing.SpanContext {
	return NewSpanContext(nil)
}

// Extract deserializes span context from from a given carrier using the format descriptor
// that tells tracer how to decode it from the carrier parameters
func (tracer *Tracer) Extract(carrier interface{}, format string) (tracing.SpanContext, error) {
	extractor, ok := tracer.extractionFormats[format]
	if !ok {
		return nil, tracing.NewUnregisteredFormatError("No extractor registered for format", format)
	}

	ctrl, ok := extractor.(Extractor)
	if !ok {
		return nil, errors.New("extractor does not have access to tracing")
	}

	ctrl.SetTracing(tracer.tracing)
	return ctrl.Extract(carrier)
}

// Inject implicitly serializes current span context using the format descriptor that
// tells how to encode trace info in the carrier parameters
func (tracer *Tracer) Inject(carrier interface{}, format string) error {
	span := tracer.currentSpan
	if span == nil {
		return nil
	}

	injector, ok := tracer.injectionFormats[format]
	if !ok {
		return tracing.NewUnregisteredFormatError("No injector registered for format", format)
	}

	return injector.Inject(span.Context(), carrier)
}

// InjectContext serializes specified span context into a given carrier using the format descriptor
// that tells how to encode trace info in the carrier parameters
func (tracer *Tracer) InjectContext(carrier interface{}, format string, spanCtx tracing.SpanContext) error {
	injector, ok := tracer.injectionFormats[format]
	if !ok {
		return tracing.NewUnregisteredFormatError("No injector registered for format", format)
	}

	return injector.Inject(spanCtx, carrier)
}

// RegisterExtractionFormat register extractor implementation for given format string
func (tracer *Tracer) RegisterExtractionFormat(format string, extractor tracing.Extractor) {
	tracer.extractionFormats[format] = extractor
}

// RegisterInjectionFormat register injector implementation for given format string
func (tracer *Tracer) RegisterInjectionFormat(format string, injector tracing.Injector) {
	tracer.injectionFormats[format] = injector
}

// Flush may flush any pending spans to the transport and reset the state of the tracer.
// Make sure this method is always called after the request is finished.
func (tracer *Tracer) Flush() {
	tracer.rootSpan = nil
	tracer.currentSpan = nil
	tracer.uuid = ""
}

// Close does a clean shutdown of the reporter, sending any traces that may be buffered in memory.
// This is especially useful for command-line tools that enable tracing,
// as well as for the long-running apps that support graceful shutdown.
//
// It goes without saying, but you cannot send anymore spans after calling Close,
// so you should only run this once during the lifecycle of the program.
func (tracer *Tracer) Close() error {
	return tracer.reporter.Close()
}

func resolveCollectorIP(host string) (string, error) {
	if net.ParseIP(host) != nil {
		return host, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "127.0.0.1", err
	}

	for _, ip := range ips {
		if ip.IsLoopback() {
			if ip.To4() != nil {
				return ip.String(), nil
			}

			continue
		}

		return ip.String(), nil
	}

	return "127.0.0.1", ErrCollectorIPNotFound
}

func registerDefaultExtractionFormats() map[string]tracing.Extractor {
	extractionFormats := make(map[string]tracing.Extractor)

	extractionFormats[formats.TextMap] = NewTextMapExtractor()
	extractionFormats[formats.HTTP] = NewHTTPExtractor()
	extractionFormats[formats.AMQP] = NewAMQPExtractor()
	extractionFormats[formats.GooglePubSub] = NewGooglePubSubExtractor()

	return extractionFormats
}

func registerDefaultInjectionFormats() map[string]tracing.Injector {
	injectionFormats := make(map[string]tracing.Injector)

	injectionFormats[formats.TextMap] = NewTextMapInjector()
	injectionFormats[formats.HTTP] = NewHTTPInjector()
	injectionFormats[formats.AMQP] = NewAMQPInjector()
	injectionFormats[formats.GooglePubSub] = NewGooglePubSubInjector()

	return injectionFormats
}
