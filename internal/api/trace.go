package api

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// TraceEnhancer is a middleware that will set the http method and url as operation-name of the top level span.
func TraceEnhancer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		resourceName := r.URL.Path
		operationName := r.Method + " " + resourceName
		span := trace.SpanFromContext(r.Context())
		if span.SpanContext().IsValid() {
			span.SetName(operationName)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
