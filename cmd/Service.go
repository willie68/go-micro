package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"wkla.no-ip.biz/go-micro/api/routes"
	"wkla.no-ip.biz/go-micro/error/serror"
	"wkla.no-ip.biz/go-micro/health"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	config "wkla.no-ip.biz/go-micro/config"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httptracer"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/go-micro/api/handler"
	"wkla.no-ip.biz/go-micro/crypt"
	clog "wkla.no-ip.biz/go-micro/logging"

	flag "github.com/spf13/pflag"
)

/*
apVersion implementing api version for this service
*/
const apiVersion = "1"
const servicename = "gomicro"

var port int
var sslport int
var system string
var serviceURL string
var registryURL string
var apikey string
var ssl bool
var configFile string
var serviceConfig config.Config
var consulAgent *consulApi.Agent
var Tracer opentracing.Tracer

func init() {
	// variables for parameter override
	ssl = false
	clog.Logger.Info("init service")
	flag.IntVarP(&port, "port", "p", 0, "port of the http server.")
	flag.IntVarP(&sslport, "sslport", "t", 0, "port of the https server.")
	flag.StringVarP(&system, "systemid", "s", "", "this is the systemid of this service. Used for the apikey generation")
	flag.StringVarP(&configFile, "config", "c", config.File, "this is the path and filename to the config file")
	flag.StringVarP(&serviceURL, "serviceURL", "u", "", "service url from outside")
	flag.StringVarP(&registryURL, "registryURL", "r", "", "registry url where to connect to consul")
}

func apiRoutes() *chi.Mux {
	myHandler := handler.NewSysAPIHandler(serviceConfig.SystemID, apikey)
	baseURL := fmt.Sprintf("/api/v%s", apiVersion)
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
		myHandler.Handler,
		httptracer.Tracer(Tracer, httptracer.Config{
			ServiceName:    servicename,
			ServiceVersion: apiVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
				//return r.URL.Path == "/healthz"
			},
			Tags: map[string]interface{}{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}),
	)

	router.Route("/", func(r chi.Router) {
		r.Mount(baseURL+"/config", routes.ConfigRoutes())
		r.Mount("/health", health.Routes())
	})
	return router
}

func healthRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
		httptracer.Tracer(Tracer, httptracer.Config{
			ServiceName:    servicename,
			ServiceVersion: apiVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
				//return r.URL.Path == "/health"
			},
			Tags: map[string]interface{}{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}),
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/", health.Routes())
	})
	return router
}

func main() {
	clog.Logger.Info("starting server")
	defer clog.Logger.Close()

	flag.Parse()

	serror.Service = servicename
	config.File = configFile
	if err := config.Load(); err != nil {
		clog.Logger.Alertf("can't load config file: %s", err.Error())
		os.Exit(1)
	}

	serviceConfig = config.Get()
	initConfig()
	initGraylog()
	var closer io.Closer
	Tracer, closer = initJaeger(servicename, serviceConfig.OpenTracing)
	opentracing.SetGlobalTracer(Tracer)
	defer closer.Close()

	healthCheckConfig := health.CheckConfig(serviceConfig.HealthCheck)

	health.InitHealthSystem(healthCheckConfig, Tracer)

	if serviceConfig.SystemID == "" {
		clog.Logger.Fatal("system id not given, can't start! Please use config file or -s parameter")
	}

	gc := crypt.GenerateCertificate{
		Organization: "EASY SOFTWARE",
		Host:         "127.0.0.1",
		ValidFor:     10 * 365 * 24 * time.Hour,
		IsCA:         false,
		EcdsaCurve:   "P384",
		Ed25519Key:   false,
	}

	if serviceConfig.Sslport > 0 {
		ssl = true
		clog.Logger.Info("ssl active")
	}

	routes.SystemID = serviceConfig.SystemID
	apikey = getApikey()
	routes.APIKey = apikey
	clog.Logger.Infof("systemid: %s", serviceConfig.SystemID)
	clog.Logger.Infof("apikey: %s", apikey)
	clog.Logger.Infof("ssl: %t", ssl)
	clog.Logger.Infof("serviceURL: %s", serviceConfig.ServiceURL)
	if serviceConfig.RegistryURL != "" {
		clog.Logger.Infof("registryURL: %s", serviceConfig.RegistryURL)
	}
	clog.Logger.Info("gomicro api routes")
	router := apiRoutes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		clog.Logger.Infof("%s %s", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		clog.Logger.Alertf("could not walk api routes. %s", err.Error())
	}
	clog.Logger.Info("health api routes")
	healthRouter := healthRoutes()
	if err := chi.Walk(healthRouter, walkFunc); err != nil {
		clog.Logger.Alertf("could not walk health routes. %s", err.Error())
	}

	var sslsrv *http.Server
	var srv *http.Server
	if ssl {
		tlsConfig, err := gc.GenerateTLSConfig()
		if err != nil {
			clog.Logger.Alertf("could not create tls config. %s", err.Error())
		}
		sslsrv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Sslport),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
			TLSConfig:    tlsConfig,
		}
		go func() {
			clog.Logger.Infof("starting https server on address: %s", sslsrv.Addr)
			if err := sslsrv.ListenAndServeTLS("", ""); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      healthRouter,
		}
		go func() {
			clog.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	} else {
		// own http server for the healthchecks
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		}
		go func() {
			clog.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	}

	if serviceConfig.RegistryURL != "" {
		initRegistry()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	clog.Logger.Info("waiting for clients")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	if ssl {
		sslsrv.Shutdown(ctx)
	}

	clog.Logger.Info("finished")

	os.Exit(0)
}

func initGraylog() {
	clog.Logger.GelfURL = serviceConfig.Logging.Gelfurl
	clog.Logger.GelfPort = serviceConfig.Logging.Gelfport
	clog.Logger.SystemID = serviceConfig.SystemID

	clog.Logger.InitGelf()
}

func initRegistry() {
	//register to consul, if configured
	consulConfig := consulApi.DefaultConfig()
	consulURL, err := url.Parse(serviceConfig.RegistryURL)
	consulConfig.Scheme = consulURL.Scheme
	consulConfig.Address = fmt.Sprintf("%s:%s", consulURL.Hostname(), consulURL.Port())
	consulClient, err := consulApi.NewClient(consulConfig)
	if err != nil {
		clog.Logger.Alertf("can't connect to consul. %v", err)
	}
	consulAgent = consulClient.Agent()

	check := new(consulApi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("%s/health/health", serviceConfig.ServiceURL)
	check.Timeout = (time.Minute * 1).String()
	check.Interval = (time.Second * 30).String()
	check.TLSSkipVerify = true
	serviceDef := &consulApi.AgentServiceRegistration{
		Name:  servicename,
		Check: check,
	}

	err = consulAgent.ServiceRegister(serviceDef)

	if err != nil {
		clog.Logger.Alertf("could not register to consul. %s", err)
		time.Sleep(time.Second * 60)
	}

}

func initConfig() {
	if port > 0 {
		serviceConfig.Port = port
	}
	if sslport > 0 {
		serviceConfig.Sslport = sslport
	}
	if system != "" {
		serviceConfig.SystemID = system
	}
	if serviceURL != "" {
		serviceConfig.ServiceURL = serviceURL
	}
}

func initJaeger(servicename string, config config.OpenTracing) (opentracing.Tracer, io.Closer) {

	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: config.Host,
			CollectorEndpoint:  config.Endpoint,
		},
	}
	if (config.Endpoint == "") && (config.Host == "") {
		cfg.Disabled = true
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func getApikey() string {
	value := fmt.Sprintf("%s_%s", servicename, serviceConfig.SystemID)
	apikey := fmt.Sprintf("%x", md5.Sum([]byte(value)))
	return strings.ToLower(apikey)
}
