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

var (
	postConfigCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gomicro_post_config_total",
		Help: "The total number of post config requests",
	})
)

// ConfigHandler the config handler
type ConfigHandler struct {
}

/*
ConfigDescription describres all metadata of a config
*/
type ConfigDescription struct {
	StoreID  string `json:"storeid"`
	TenantID string `json:"tenantID"`
	Size     int    `json:"size"`
}

// NewConfigHandler creates a new REST config handler
func NewConfigHandler() api.Handler {
	return &ConfigHandler{}
}

// Routes getting all routes for the config endpoint
func (c *ConfigHandler) Routes() (string, *chi.Mux) {
	router := chi.NewRouter()
	router.Post("/", c.PostConfig)
	router.Get("/", c.GetConfig)
	router.Delete("/", c.DeleteConfig)
	router.Get("/size", c.GetConfigSize)
	return BaseURL + configSubpath, router
}

// GetConfig getting if a store for a tenant is initialized
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
func (c *ConfigHandler) GetConfig(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	cd := ConfigDescription{
		StoreID:  "myNewStore",
		TenantID: tenant,
		Size:     1234567,
	}
	render.JSON(response, request, cd)
}

// PostConfig create a new store for a tenant
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
func (c *ConfigHandler) PostConfig(response http.ResponseWriter, request *http.Request) {
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

// DeleteConfig deleting store for a tenant, this will automatically delete all data in the store
// @Summary Delete a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 "ok"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Router /config [delete]
func (c *ConfigHandler) DeleteConfig(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	render.JSON(response, request, tenant)
}

// GetConfigSize size of the store for a tenant
// @Summary Get size of a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 {string} string "size"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Router /config/size [get]
func (c *ConfigHandler) GetConfigSize(response http.ResponseWriter, request *http.Request) {
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
