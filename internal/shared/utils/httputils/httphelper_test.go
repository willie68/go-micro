package httputils

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/adapter/inbound/http/api"
	"github.com/willie68/go-micro/internal/adapter/inbound/http/auth"
	"github.com/willie68/go-micro/internal/shared/serror"
)

type sampleBody struct {
	Name string `json:"name" validate:"required"`
}

func TestTenantIDFromHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := TenantID(req)
	assert.Error(t, err)

	req.Header.Set(api.TenantHeaderKey, "Acme")
	id, err := TenantID(req)
	assert.NoError(t, err)
	assert.Equal(t, "acme", id)
}

func TestTenantIDFromURLParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(api.URLParamTenantID, "TenantA")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	id, err := TenantID(req)
	assert.NoError(t, err)
	assert.Equal(t, "tenanta", id)
}

func TestTenantIDFromJWT(t *testing.T) {
	TenantClaim = "tenant"
	Strict = false
	jt := &auth.JWT{Payload: map[string]any{"tenant": "JwtTenant"}}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(auth.NewContext(req.Context(), jt, nil, false))
	id, err := TenantID(req)
	assert.NoError(t, err)
	assert.Equal(t, "jwttenant", id)

	Strict = true
	jt = &auth.JWT{Payload: map[string]any{}}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(auth.NewContext(req.Context(), jt, nil, false))
	_, err = TenantID(req)
	assert.Error(t, err)
	Strict = false
}

func TestParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := Param(req, "id")
	assert.Error(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	id, err := Param(req, "id")
	assert.NoError(t, err)
	assert.Equal(t, "42", id)
}

func TestDecode(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"ok"}`))
	req.Header.Set("Content-Type", "application/json")
	var body sampleBody
	assert.NoError(t, Decode(req, &body))
	assert.Equal(t, "ok", body.Name)

	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	assert.Error(t, Decode(req, &sampleBody{}))
}

func TestCreatedAndErr(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	Created(rec, req, "abc", map[string]string{"id": "abc"})
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "abc", rec.Header().Get("Location"))

	rec = httptest.NewRecorder()
	Err(rec, req, serror.NotFound("address", "1"))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFileServer(t *testing.T) {
	r := chi.NewRouter()
	assert.Panics(t, func() {
		FileServer(r, "/{id}", http.Dir("."))
	})
	FileServer(r, "/static", http.Dir("."))
}
