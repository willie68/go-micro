package apiv1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/serror"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/pkg/pmodel"

	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/utils/httputils"
)

const (
	errMissingTenantKey = "missing-tenant"
	errMissingTenantMsg = "tenant header %s missing"
)

var (
	postAdrCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gomicro_post_adr_total",
		Help: "The total number of address requests",
	})
)

// AdrHandler the address handler
type AdrHandler struct {
	adrstg adrsvc.AddressStorage
}

// NewAdrHandler creates a new REST address handler
func NewAdrHandler() api.Handler {
	return &AdrHandler{
		adrstg: do.MustInvoke[adrsvc.AddressStorage](nil),
	}
}

// Routes getting all routes for the address endpoint
func (c *AdrHandler) Routes() (string, *chi.Mux) {
	router := chi.NewRouter()
	router.Post("/", c.PostAddress)
	router.Get("/", c.GetAddresses)
	router.Get("/{id}", c.GetAddress)
	router.Post("/{id}", c.UpdateAddress)
	router.Delete("/{id}", c.DeleteAddress)
	return BaseURL + addressesSubpath, router
}

// GetAddresses getting all addresses
//
//	@Summary	getting all addresses
//	@Tags		addresses
//	@Accept		json
//	@Produce	json
//	@Security	api_key
//	@Param		tenant	header		string			true	"Tenant"
//	@Success	200		{array}		pmodel.Address	"response with list of addresses as json"
//	@Failure	400		{object}	serror.Serr		"client error information as json"
//	@Failure	500		{object}	serror.Serr		"server error information as json"
//	@Router		/addresses [get]
func (c *AdrHandler) GetAddresses(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf(errMissingTenantMsg, api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, errMissingTenantKey, msg))
		return
	}
	l, err := c.adrstg.Addresses()
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	render.JSON(response, request, l)
}

// GetAddress getting one address
//
//	@Summary	getting one address
//	@Tags		addresses
//	@Accept		json
//	@Produce	json
//	@Security	api_key
//	@Param		tenant	header		string			true	"Tenant"
//	@Param		id		path		string			true	"ID"
//	@Success	200		{array}		pmodel.Address	"response with the address with id as json"
//	@Failure	400		{object}	serror.Serr		"client error information as json"
//	@Failure	500		{object}	serror.Serr		"server error information as json"
//	@Router		/addresses/{id} [get]
func (c *AdrHandler) GetAddress(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf(errMissingTenantMsg, api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, errMissingTenantKey, msg))
		return
	}
	n := chi.URLParam(request, "id")
	if !c.adrstg.Has(n) {
		httputils.Err(response, request, serror.ErrNotExists)
		return
	}
	adr, err := c.adrstg.Read(n)
	if err != nil {
		httputils.Err(response, request, serror.ErrUnknowError)
		return
	}

	render.JSON(response, request, adr)
}

// PostAddress create a new address, this method will always return 201
//
//	@Summary	Create a new address
//	@Tags		addresses
//	@Accept		json
//	@Produce	json
//	@Security	api_key
//	@Param		tenant	header		string			true	"Tenant"
//	@Param		payload	body		pmodel.Address	true	"address to be added"
//	@Success	201		{string}	string			"tenant"
//	@Failure	400		{object}	serror.Serr		"client error information as json"
//	@Failure	500		{object}	serror.Serr		"server error information as json"
//	@Router		/addresses [post]
func (c *AdrHandler) PostAddress(response http.ResponseWriter, request *http.Request) {
	var b []byte
	var err error
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf(errMissingTenantMsg, api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, errMissingTenantKey, msg))
		return
	}
	logger.Infof("create config: tenant %s", tenant)

	if b, err = io.ReadAll(request.Body); err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	var adr pmodel.Address

	err = json.Unmarshal(b, &adr)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	postAdrCounter.Inc()
	n, err := c.adrstg.Create(adr)
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

// UpdateAddress update an address, this method will always return 201
//
//	@Summary	Create a new address
//	@Tags		addresses
//	@Accept		json
//	@Produce	json
//	@Security	api_key
//	@Param		tenant	header		string		true	"Tenant"
//	@Param		payload	body		string		true	"Add store"
//	@Success	201		{string}	string		"tenant"
//	@Failure	400		{object}	serror.Serr	"client error information as json"
//	@Failure	500		{object}	serror.Serr	"server error information as json"
//	@Router		/addresses/{id} [post]
func (c *AdrHandler) UpdateAddress(response http.ResponseWriter, request *http.Request) {
	var b []byte
	var err error
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf(errMissingTenantMsg, api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, errMissingTenantKey, msg))
		return
	}
	logger.Infof("create config: tenant %s", tenant)

	if b, err = io.ReadAll(request.Body); err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	var adr pmodel.Address

	err = json.Unmarshal(b, &adr)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusBadRequest))
		return
	}
	postAdrCounter.Inc()
	err = c.adrstg.Update(adr)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	render.Status(request, http.StatusCreated)
	render.JSON(response, request, adr)
}

// DeleteAddress deleting address
//
//	@Summary	Delete a address
//	@Tags		addresses
//	@Accept		json
//	@Produce	json
//	@Security	api_key
//	@Param		tenant	header	string	true	"Tenant"
//	@Success	200		"ok"
//	@Failure	400		{object}	serror.Serr	"client error information as json"
//	@Router		/addresses/{id} [delete]
func (c *AdrHandler) DeleteAddress(response http.ResponseWriter, request *http.Request) {
	tenant := getTenant(request)
	if tenant == "" {
		msg := fmt.Sprintf(errMissingTenantMsg, api.TenantHeaderKey)
		httputils.Err(response, request, serror.BadRequest(nil, errMissingTenantKey, msg))
		return
	}
	n := chi.URLParam(request, "id")
	if !c.adrstg.Has(n) {
		httputils.Err(response, request, serror.NotFound("address", n))
		return
	}
	adr, err := c.adrstg.Read(n)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	err = c.adrstg.Delete(n)
	if err != nil {
		httputils.Err(response, request, serror.Wrapc(err, http.StatusInternalServerError))
		return
	}
	render.JSON(response, request, adr)
}

// getTenant getting the tenant from the request
func getTenant(req *http.Request) string {
	return req.Header.Get(api.TenantHeaderKey)
}
