package apiv1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httptracer"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/auth"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/utils/httputils"
	"github.com/willie68/go-micro/pkg/web"
)

// APIVersion the actual implemented api version
const APIVersion = "1"

// BaseURL is the url all endpoints will be available under
var BaseURL = fmt.Sprintf("/api/v%s", APIVersion)

var logger = logging.New().WithName("apiv1")

// defining all sub pathes for api v1
const addressesSubpath = "/addresses"

func token(r *http.Request) (string, error) {
	tk := r.Header.Get("Authorization")
	tk = strings.TrimPrefix(tk, "Bearer ")
	return tk, nil
}

// APIRoutes configuring the api routes for the main REST API
func APIRoutes(cfn config.Config, trc opentracing.Tracer) (*chi.Mux, error) {
	logger.Infof("baseurl : %s", BaseURL)
	router := chi.NewRouter()
	setDefaultHandler(router, cfn, trc)

	// jwt is activated, register the Authenticator and Validator
	if strings.EqualFold(cfn.Auth.Type, "jwt") {
		err := setJWTHandler(router, cfn)
		if err != nil {
			return nil, err
		}
	}

	// building the routes
	router.Route("/", func(r chi.Router) {
		r.Mount(NewAdrHandler().Routes())

		r.Mount(health.NewHealthHandler().Routes())
		if cfn.Metrics.Enable {
			r.Mount("/metrics", promhttp.Handler())
		}
	})
	// adding a file server with web client asserts
	httputils.FileServer(router, "/client", http.FS(web.WebClientAssets))
	logger.Infof("%s api routes", config.Servicename)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Infof("api route: %s %s", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Alertf("could not walk api routes. %s", err.Error())
	}
	return router, nil
}

func setJWTHandler(router *chi.Mux, cfn config.Config) error {
	jwtConfig, err := auth.ParseJWTConfig(cfn.Auth)
	if err != nil {
		return err
	}
	logger.Infof("jwt config: %v", jwtConfig)
	jwtAuth := auth.JWTAuth{
		Config: jwtConfig,
	}
	router.Use(
		auth.Verifier(&jwtAuth),
		auth.Authenticator,
	)
	return nil
}

func setDefaultHandler(router *chi.Mux, cfn config.Config, tracer opentracing.Tracer) {
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Recoverer,
		cors.Handler(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-mcs-username", "X-mcs-password", "X-mcs-profile"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	if tracer != nil {
		router.Use(httptracer.Tracer(tracer, httptracer.Config{
			ServiceName:    config.Servicename,
			ServiceVersion: "V" + APIVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
				//return r.URL.Path == "/livez"
			},
			Tags: map[string]any{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}))
	}
	if cfn.Metrics.Enable {
		router.Use(
			api.MetricsHandler(api.MetricsConfig{
				SkipFunc: func(r *http.Request) bool {
					return false
				},
			}),
		)
	}
}

// HealthRoutes returning the health routes
func HealthRoutes(cfn config.Config, tracer opentracing.Tracer) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Recoverer,
	)
	if tracer != nil {
		router.Use(httptracer.Tracer(tracer, httptracer.Config{
			ServiceName:    config.Servicename,
			ServiceVersion: "V" + APIVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
			},
			Tags: map[string]any{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}))
	}
	if cfn.Metrics.Enable {
		router.Use(
			api.MetricsHandler(api.MetricsConfig{
				SkipFunc: func(r *http.Request) bool {
					return false
				},
			}),
		)
	}

	router.Route("/", func(r chi.Router) {
		r.Mount(health.NewHealthHandler().Routes())
		if cfn.Metrics.Enable {
			r.Mount(api.MetricsEndpoint, promhttp.Handler())
		}
		if cfn.Profiling.Enable {
			// Define the routes for serving profiling data
			r.Mount("/debug", middleware.Profiler())
		}
	})

	logger.Info("health api routes")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Infof("health route: %s %s", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Alertf("could not walk health routes. %s", err.Error())
	}

	return router
}
