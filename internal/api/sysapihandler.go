package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/willie68/go-micro/internal/serror"
)

// SysAPIConfig defining a handler for checking system id and api key
type SysAPIConfig struct {
	Apikey string
	// Skip particular requests from the handler
	SkipFunc func(r *http.Request) bool
}

// SysAPIHandler creates a new directly usable handler
func SysAPIHandler(cfg SysAPIConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip tracer
			if cfg.SkipFunc != nil && cfg.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			if cfg.Apikey != strings.ToLower(r.Header.Get(APIKeyHeaderKey)) {
				msg := "apikey not correct"
				apierr := serror.BadRequest(nil, "missing-header", msg)
				render.Status(r, apierr.Code)
				render.JSON(w, r, apierr)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
