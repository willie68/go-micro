package api

import "github.com/go-chi/chi/v5"

// TenantHeaderKey in this header the right tenant should be inserted
const TenantHeaderKey = "tenant"

// URLParamTenantID url parameter for the tenant id
const URLParamTenantID = "tntid"

// MetricsEndpoint endpoint subpath  for metrics
const MetricsEndpoint = "/metrics"

// Handler a http REST interface handler
type Handler interface {
	// Routes get the routes
	Routes() (string, *chi.Mux)
}
