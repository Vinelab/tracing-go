package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Vinelab/tracing-go"
	"github.com/Vinelab/tracing-go/formats"
	"github.com/Vinelab/tracing-go/support/slice"
	"github.com/go-chi/chi/middleware"
)

// TraceRequests middleware
type TraceRequests struct {
	tracer       tracing.Tracer
	contentTypes []string
	excludedPaths []string
}

// NewTraceRequests creates a new TraceRequests middleware with the provided options
func NewTraceRequests(tracer tracing.Tracer, contentTypes []string, excludedPaths []string) *TraceRequests {
	return &TraceRequests{
		tracer:       tracer,
		contentTypes: contentTypes,
		excludedPaths: excludedPaths,
	}
}

// Handler applies tracing on the request and ensures that we collect metadata from
// a request-response cycle. This includes uri and method of the request, headers,
// client ip, input, response code and content etc.
func (mdlw *TraceRequests) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if slice.Contains(mdlw.excludedPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Create a proxy that hooks into response and allows us to access its contents
		response := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Save response body in the buffer for logging purposes
		buffer := bytes.Buffer{}
		response.Tee(&buffer)

		// Extract existing trace from request headers (if present)
		spanContext, err := mdlw.tracer.Extract(r, formats.HTTP)
		if err != nil {
			log.Fatal(err)
		}

		// Start the global span, it'll wrap the request lifecycle.
		// span := tr.tracer.StartSpan("HTTP Request", spanContext)
		span := mdlw.tracer.StartSpan("HTTP Request", spanContext)

		// Save request metadata for this span. Note that tags are searchable on UI.
		span.Tag("type", "http")
		span.Tag("request_method", r.Method)
		span.Tag("request_path", r.URL.Path)
		span.Tag("request_uri", r.RequestURI)
		span.Tag("request_headers", getHeaders(r.Header))
		span.Tag("request_ip", strings.Split(r.RemoteAddr, ":")[0])
		if slice.Contains(mdlw.contentTypes, r.Header.Get("Content-Type")) {
			span.Tag("request_input", getRequestInput(r))
		}

		defer func() {
			span.Tag("response_status", strconv.Itoa(response.Status()))
			span.Tag("response_headers", getHeaders(response.Header()))

			if slice.Contains(mdlw.contentTypes, response.Header().Get("Content-Type")) {
				span.Tag("response_content", buffer.String())
			}

			span.Finish()
			mdlw.tracer.Flush()
		}()

		next.ServeHTTP(response, r)
	}

	return http.HandlerFunc(fn)
}

func getRequestInput(r *http.Request) string {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Unable to read request body")
	}

	return string(data)
}

func getHeaders(h http.Header) string {
	str := ""
	for key, value := range h {
		str = fmt.Sprintln(str, key+": "+value[0])
	}
	return str
}
