package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"wkla.no-ip.biz/go-micro/error/serror"
)

// Validate validator
var Validate *validator.Validate

// TokenHeader in this header the token is expected
const TokenHeader = "Authorization"

// TenantHeader in this header the tenant is expected
const TenantHeader = "X-es-tenant"

// UserHeader in this header the username is expected
const UserHeader = "X-es-username"

// CustomerHeader in this header the customer is expected (only for tenant unspecific requests)
const CustomerHeader = "X-es-customer"

// Username gets the username of the given request
func Username(r *http.Request) (string, error) {
	uid := r.Header.Get(UserHeader)
	if uid == "" {
		msg := fmt.Sprintf("user header missing: %s", UserHeader)
		return "", serror.BadRequest(nil, "missing-header", msg)
	}
	return uid, nil
}

// TenantID gets the tenant-id of the given request
func TenantID(r *http.Request) (string, error) {
	id := r.Header.Get(TenantHeader)
	if id == "" {
		msg := fmt.Sprintf("tenant header %s missing", TenantHeader)
		return "", serror.BadRequest(nil, "missing-tenant", msg)
	}
	return id, nil
}

// CustomerID gets the customer-id of the given request
func CustomerID(r *http.Request) (string, error) {
	id := r.Header.Get(CustomerHeader)
	if id == "" {
		msg := fmt.Sprintf("customer header %s missing", TenantHeader)
		return "", serror.BadRequest(nil, "missing-customer", msg)
	}
	return id, nil
}

// Decode decodes and validates an object
func Decode(r *http.Request, v interface{}) error {
	err := render.DefaultDecoder(r, v)
	if err != nil {
		serror.BadRequest(err, "decode-body", "could not decode body")
	}
	if err := Validate.Struct(v); err != nil {
		serror.BadRequest(err, "validate-body", "body invalid")
	}
	return nil
}

// Param gets the url param of the given request
func Param(r *http.Request, name string) (string, error) {
	cid := chi.URLParam(r, name)
	if cid == "" {
		msg := fmt.Sprintf("missing %s in path", name)
		return "", serror.BadRequest(nil, "missing-param", msg)
	}
	return cid, nil
}

// Created object created
func Created(w http.ResponseWriter, r *http.Request, id string, v interface{}) {
	// TODO add relative path to location
	w.Header().Add("Location", fmt.Sprintf("%s", id))
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, v)
}

// Err writes an error response
func Err(w http.ResponseWriter, r *http.Request, err error) {
	apierr := serror.Wrap(err, "unexpected-error")
	render.Status(r, apierr.Code)
	render.JSON(w, r, apierr)
}

func init() {
	Validate = validator.New()
}
