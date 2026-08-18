package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/willie68/go-micro/internal/services/logging"
)

// TenantHeaderKey in this header the right tenant should be inserted
const TenantHeaderKey = "tenant"

// URLParamTenantID url parameter for the tenant id
const URLParamTenantID = "tntid"

// MetricsEndpoint endpoint subpath  for metrics
const MetricsEndpoint = "/metrics"

var logger = logging.New("api")

// Handler a http REST interface handler
type Handler interface {
	// Routes get the routes
	Routes() (string, *chi.Mux)
}
