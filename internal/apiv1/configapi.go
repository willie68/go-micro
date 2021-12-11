package apiv1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/willie68/go-micro/internal/serror"

	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/utils/httputils"
)

// TenantHeader in this header thr right tenant should be inserted
const timeout = 1 * time.Minute

//APIKey the apikey of this service
var APIKey string

//SystemID the systemid of this service
var SystemID string

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

/*
GetConfigEndpoint getting if a store for a tenant is initialised
because of the automatic store creation, the value is more likely that data is stored for this tenant
*/
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

/*
PostConfigEndpoint create a new store for a tenant
because of the automatic store creation, this method will always return 201
*/
func PostConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	log.Printf("create store for tenant %s", tenant)
	render.Status(request, http.StatusCreated)
	render.JSON(response, request, tenant)
}

/*
DeleteConfigEndpoint deleting store for a tenant, this will automatically delete all data in the store
*/
func DeleteConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, "missing-tenant", msg))
		return
	}
	render.JSON(response, request, tenant)
}

/*
GetConfigSizeEndpoint size of the store for a tenant
*/
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
