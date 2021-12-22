package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsConfig defining a handler for checking system id and api key
type MetricsConfig struct {
	// Skip particular requests from the handler
	SkipFunc func(r *http.Request) bool
}

var(
	metrics map[string]prometheus.Counter = make(map[string]prometheus.Counter)
	m = regexp.MustCompile("[^a-zA-Z_]")
) 

// MetricsHandler creates a new directly usable handler
func MetricsHandler(cfg MetricsConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip tracer
			if cfg.SkipFunc != nil && cfg.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}
			resourceName := r.URL.Path
			resourceName = strings.TrimSuffix(resourceName, "/")
			resourceName = strings.TrimPrefix(resourceName, "/")
			resourceName = m.ReplaceAllString(resourceName, "_")
			metricsName := fmt.Sprintf("rest_requests_%s_%s_total", r.Method, resourceName)
			counter, ok := metrics[metricsName]
			if !ok {
				counter = promauto.NewCounter(prometheus.CounterOpts{
					Name: metricsName,
					Help: fmt.Sprintf("auto generated metrics for \"%s\" with method \"%s\"", r.URL.Path, r.Method),
				})
				metrics[metricsName] = counter
			}
			counter.Inc()
			next.ServeHTTP(w, r)
		})
	}
}
