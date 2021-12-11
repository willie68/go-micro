package api

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// TraceEnhancer is a middleware that will set the http method and url as operation-name of the top level span.
func TraceEnhancer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		resourceName := r.URL.Path
		operationName := r.Method + " " + resourceName
		span := opentracing.SpanFromContext(r.Context())
		if span != nil {
			span.SetOperationName(operationName)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
