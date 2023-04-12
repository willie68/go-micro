package apiv1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/model"
	"github.com/willie68/go-micro/internal/serror"
	"github.com/willie68/go-micro/internal/services/sconfig"
	"github.com/willie68/go-micro/pkg/pmodel"

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
	cfgs sconfig.SConfig
}

// NewConfigHandler creates a new REST config handler
func NewConfigHandler() api.Handler {
	return &ConfigHandler{
		cfgs: do.MustInvokeNamed[sconfig.SConfig](nil, sconfig.DoConfig),
	}
}

// Routes getting all routes for the config endpoint
func (c *ConfigHandler) Routes() (string, *chi.Mux) {
	router := chi.NewRouter()
	router.Post("/", c.PostConfig)
	router.Get("/", c.GetConfigs)
	router.Get("/{id}", c.GetConfig)
	router.Delete("/{id}", c.DeleteConfig)
	router.Get("/_own", c.GetConfigOfTenant)
	return BaseURL + configSubpath, router
}

// GetConfigs getting all configs
// @Summary getting all configs
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 {array} ConfigDescription "response with config as json"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /config [get]
func (c *ConfigHandler) GetConfigs(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	l, err := c.cfgs.List()
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	render.JSON(response, request, l)
}

// GetConfig getting one configs
// @Summary getting one configs
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
	n := chi.URLParam(request, "id")
	if !c.cfgs.HasConfig(n) {
		httputils.Err(response, request, serror.ErrNotExists)
		return
	}
	cd, err := c.cfgs.GetConfig(n)
	if err != nil {
		httputils.Err(response, request, serror.ErrUnknowError)
		return
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
	var b []byte
	var err error
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	log.Printf("create config: tenant %s", tenant)

	if b, err = io.ReadAll(request.Body); err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	var cd pmodel.ConfigDescription

	err = json.Unmarshal(b, &cd)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	cdm := model.ConfigDescription{
		StoreID:  cd.StoreID,
		TenantID: cd.TenantID,
		Size:     cd.Size,
	}
	postConfigCounter.Inc()
	n, err := c.cfgs.PutConfig(cd.TenantID, cdm)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	id := struct {
		ID string `json:"id"`
	}{
		ID: n,
	}

	render.Status(request, http.StatusCreated)
	render.JSON(response, request, id)
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
	n := chi.URLParam(request, "id")
	if !c.cfgs.HasConfig(n) {
		httputils.Err(response, request, serror.NotFound("config", n))
		return
	}
	c.cfgs.DeleteConfig(n)
	render.JSON(response, request, tenant)
}

// GetConfigOfTenant config of tenant
// @Summary Get size of a store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Success 200 {string} string "size"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Router /config/size [get]
func (c *ConfigHandler) GetConfigOfTenant(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	cd, err := c.cfgs.GetConfig(tenant)
	if err != nil {
		httputils.Err(response, request, serror.ErrUnknowError)
		return
	}
	render.JSON(response, request, cd)
}

// getTenant getting the tenant from the request
func getTenant(req *http.Request) string {
	return req.Header.Get(api.TenantHeaderKey)
}
