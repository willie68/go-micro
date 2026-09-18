package apiv1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/adapter/inbound/http/api"
	"github.com/willie68/go-micro/internal/infrastructure/logging"
	"github.com/willie68/go-micro/pkg/pmodel"
)

type fakeAddress struct {
	items map[string]pmodel.Address
	err   error
}

func (f *fakeAddress) Addresses(context.Context) ([]pmodel.Address, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := make([]pmodel.Address, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, v)
	}
	return out, nil
}

func (f *fakeAddress) Has(_ context.Context, id string) bool {
	_, ok := f.items[id]
	return ok
}

func (f *fakeAddress) Read(_ context.Context, id string) (*pmodel.Address, error) {
	if f.err != nil {
		return nil, f.err
	}
	adr, ok := f.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &adr, nil
}

func (f *fakeAddress) Create(_ context.Context, adr pmodel.Address) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	adr.ID = "created-1"
	f.items[adr.ID] = adr
	return adr.ID, nil
}

func (f *fakeAddress) Update(_ context.Context, adr pmodel.Address) error {
	if f.err != nil {
		return f.err
	}
	f.items[adr.ID] = adr
	return nil
}

func (f *fakeAddress) Delete(_ context.Context, id string) error {
	if f.err != nil {
		return f.err
	}
	delete(f.items, id)
	return nil
}

func newHandler(stg *fakeAddress) *AdrHandler {
	return &AdrHandler{
		adrstg: stg,
		logger: logging.New("test"),
	}
}

func serve(h *AdrHandler, method, path string, body []byte, tenant string) *httptest.ResponseRecorder {
	_, mux := h.Routes()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if tenant != "" {
		req.Header.Set(api.TenantHeaderKey, tenant)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestNewAdrHandler(t *testing.T) {
	inj := do.New()
	do.ProvideValue(inj, &fakeAddress{items: map[string]pmodel.Address{}})
	h := NewAdrHandler(inj)
	assert.NotNil(t, h)
	path, mux := h.Routes()
	assert.Equal(t, BaseURL+addressesSubpath, path)
	assert.NotNil(t, mux)
}

func TestGetAddresses(t *testing.T) {
	h := newHandler(&fakeAddress{items: map[string]pmodel.Address{
		"1": {ID: "1", Name: "Doe"},
	}})

	rec := serve(h, http.MethodGet, "/", nil, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodGet, "/", nil, "t1")
	assert.Equal(t, http.StatusOK, rec.Code)

	h.adrstg = &fakeAddress{items: map[string]pmodel.Address{}, err: errors.New("boom")}
	rec = serve(h, http.MethodGet, "/", nil, "t1")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAddress(t *testing.T) {
	h := newHandler(&fakeAddress{items: map[string]pmodel.Address{
		"1": {ID: "1", Name: "Doe"},
	}})

	rec := serve(h, http.MethodGet, "/1", nil, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodGet, "/missing", nil, "t1")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	rec = serve(h, http.MethodGet, "/1", nil, "t1")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPostAddress(t *testing.T) {
	h := newHandler(&fakeAddress{items: map[string]pmodel.Address{}})

	rec := serve(h, http.MethodPost, "/", []byte(`{"name":"Doe"}`), "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodPost, "/", []byte(`not-json`), "t1")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodPost, "/", []byte(`{"name":"Doe"}`), "t1")
	assert.Equal(t, http.StatusCreated, rec.Code)

	h.adrstg = &fakeAddress{items: map[string]pmodel.Address{}, err: errors.New("boom")}
	rec = serve(h, http.MethodPost, "/", []byte(`{"name":"Doe"}`), "t1")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateAddress(t *testing.T) {
	h := newHandler(&fakeAddress{items: map[string]pmodel.Address{
		"1": {ID: "1", Name: "Doe"},
	}})

	body, _ := json.Marshal(pmodel.Address{ID: "1", Name: "Updated"})
	rec := serve(h, http.MethodPost, "/1", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodPost, "/1", []byte(`not-json`), "t1")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodPost, "/1", body, "t1")
	assert.Equal(t, http.StatusCreated, rec.Code)

	h.adrstg = &fakeAddress{items: map[string]pmodel.Address{}, err: errors.New("boom")}
	rec = serve(h, http.MethodPost, "/1", body, "t1")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteAddress(t *testing.T) {
	h := newHandler(&fakeAddress{items: map[string]pmodel.Address{
		"1": {ID: "1", Name: "Doe"},
	}})

	rec := serve(h, http.MethodDelete, "/1", nil, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = serve(h, http.MethodDelete, "/missing", nil, "t1")
	assert.Equal(t, http.StatusNotFound, rec.Code)

	rec = serve(h, http.MethodDelete, "/1", nil, "t1")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.Empty(t, getTenant(req))
	req.Header.Set(api.TenantHeaderKey, "acme")
	assert.Equal(t, "acme", getTenant(req))
}
