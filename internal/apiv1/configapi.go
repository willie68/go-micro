package apiv1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willie68/go-micro/internal/serror"

	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/utils/httputils"
)

//APIKey the apikey of this service
var APIKey string

var (
	postConfigCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gomicro_post_config_total",
		Help: "The total number of post config requests",
	})
)

/*
ConfigDescription describres all metadata of a config
*/
type ConfigDescription struct {
	StoreID  string `json:"storeid"`
	TenantID string `json:"tenantID"`
	Size     int    `json:"size"`
}

/*
ConfigRoutes getting all routes for the config endpoint
*/
func ConfigRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", PostConfigEndpoint)
	router.Get("/", GetConfigEndpoint)
	router.Delete("/", DeleteConfigEndpoint)
	router.Get("/size", GetConfigSizeEndpoint)
	return router
}

// GetConfigEndpoint getting if a store for a tenant is initialised
// because of the automatic store creation, the value is more likely that data is stored for this tenant
// @Summary Get a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 {array} ConfigDescription "response with config as json"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /config [get]
func GetConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	c := ConfigDescription{
		StoreID:  "myNewStore",
		TenantID: tenant,
		Size:     1234567,
	}
	render.JSON(response, request, c)
}

// PostConfigEndpoint create a new store for a tenant
// because of the automatic store creation, this method will always return 201
// @Summary Create a new store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Param payload body string true "Add store"
// @Success 201 {string} string "tenant"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /config [post]
func PostConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	log.Printf("create store for tenant %s", tenant)
	postConfigCounter.Inc()
	render.Status(request, http.StatusCreated)
	render.JSON(response, request, tenant)
}

// DeleteConfigEndpoint deleting store for a tenant, this will automatically delete all data in the store
// @Summary Delete a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 "ok"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Router /config [delete]
func DeleteConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	render.JSON(response, request, tenant)
}

// GetConfigSizeEndpoint size of the store for a tenant
// @Summary Get size of a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 {string} string "size"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Router /config/size [get]
func GetConfigSizeEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}

	render.JSON(response, request, tenant)
}

/*
getTenant getting the tenant from the request
*/
func getTenant(req *http.Request) string {
	return req.Header.Get(api.TenantHeaderKey)
}
