package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"wkla.no-ip.biz/go-micro/error/serror"
)

// KeyHeader in this header the right api key should be inserted
const KeyHeader = "X-es-apikey"

// SystemHeader in this header the right system should be inserted
const SystemHeader = "X-es-system"

// SysAPIKey defining a handler for checking system id and api key
type SysAPIKey struct {
	SystemID string
	Apikey   string
	log      *log.Logger
}

// NewSysAPIHandler creates a new SysApikeyHandler
func NewSysAPIHandler(systemID string, apikey string) *SysAPIKey {
	return &SysAPIKey{
		SystemID: systemID,
		Apikey:   apikey,
	}
}

// Handler the handler checks systemid and apikey headers
func (s *SysAPIKey) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.URL.Path, "/")
		if !strings.HasPrefix(path, "/health") {
			if s.SystemID != r.Header.Get(SystemHeader) {
				msg := "either system id or apikey not correct"
				apierr := serror.BadRequest(nil, "missing-header", msg)
				render.Status(r, apierr.Code)
				render.JSON(w, r, apierr)
				return
			}
			if s.Apikey != strings.ToLower(r.Header.Get(KeyHeader)) {
				msg := "either system id or apikey not correct"
				apierr := serror.BadRequest(nil, "missing-header", msg)
				render.Status(r, apierr.Code)
				render.JSON(w, r, apierr)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

var (
	ContextKeyOffset = contextKey("offset")
	ContextKeyLimit  = contextKey("limit")
)

//Paginate is a middleware logic for populating the context with offset and limit values
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				msg := "type of offset string is not correct."
				apierr := serror.BadRequest(err, "wrong-type", msg)
				render.Status(request, apierr.Code)
				render.JSON(response, request, apierr)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyOffset, offset)
		} else {
			ctx = context.WithValue(ctx, ContextKeyOffset, 0)
		}
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				msg := "type of limit string is not correct."
				apierr := serror.BadRequest(err, "wrong-type", msg)
				render.Status(request, apierr.Code)
				render.JSON(response, request, apierr)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyLimit, limit)
		} else {
			ctx = context.WithValue(ctx, ContextKeyLimit, 0)
		}
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

type contextKey string

func (c contextKey) String() string {
	return "api" + string(c)
}
