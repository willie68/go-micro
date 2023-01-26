package httputils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/auth"
	"github.com/willie68/go-micro/internal/serror"
)

// val validator
var val *validator.Validate

// TenantClaim the claim used in the jwt claims, where the actual tenant is stored
var TenantClaim string

// Strict throwing an error if the tenant is not present in the token
var Strict bool

// TenantID gets the tenant-id of the given request
func TenantID(r *http.Request) (string, error) {
	tntID := chi.URLParam(r, api.URLParamTenantID)
	if tntID != "" {
		return strings.ToLower(tntID), nil
	}
	var id string
	_, claims, _ := auth.FromContext(r.Context())
	if claims != nil {
		tenant, ok := claims[TenantClaim].(string)
		if ok {
			return strings.ToLower(tenant), nil
		}
		if Strict {
			return "", serror.BadRequest(nil, "missing-tenant", "no tenant claim in jwt token")
		}
	}
	id = r.Header.Get(api.TenantHeaderKey)
	if id == "" {
		msg := fmt.Sprintf("tenant header %s missing", api.TenantHeaderKey)
		return "", serror.BadRequest(nil, "missing-tenant", msg)
	}
	return strings.ToLower(id), nil
}

// Decode decodes and validates an object
func Decode(r *http.Request, v any) error {
	err := render.DefaultDecoder(r, v)
	if err != nil {
		return serror.BadRequest(err, "decode-body", "could not decode body")
	}
	if err := val.Struct(v); err != nil {
		return serror.BadRequest(err, "validate-body", "body invalid")
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
func Created(w http.ResponseWriter, r *http.Request, id string, v any) {
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
	val = validator.New()
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		//rctx := chi.RouteContext(r.Context())
		//pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.FileServer(root)
		fs.ServeHTTP(w, r)
	})
}
