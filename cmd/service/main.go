// Package main this is the entry point into the service
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/apiv1"
	"github.com/willie68/go-micro/internal/serror"
	"github.com/willie68/go-micro/internal/services"
	"github.com/willie68/go-micro/internal/services/shttp"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	config "github.com/willie68/go-micro/internal/config"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	log "github.com/willie68/go-micro/internal/logging"

	flag "github.com/spf13/pflag"
)

var (
	port          int
	sslport       int
	serviceURL    string
	configFile    string
	serviceConfig config.Config
	tracer        opentracing.Tracer
)

func init() {
	// variables for parameter override
	log.Root.Info("init service")
	flag.IntVarP(&port, "port", "p", 0, "port of the http server.")
	flag.IntVarP(&sslport, "sslport", "t", 0, "port of the https server.")
	flag.StringVarP(&configFile, "config", "c", config.File, "this is the path and filename to the config file")
	flag.StringVarP(&serviceURL, "serviceURL", "u", "", "service url from outside")
}

// @title GoMicro service API
// @version 1.0
// @description The GoMicro service is a template for microservices written in go.
// @BasePath /api/v1
// @securityDefinitions.apikey api_key
// @in header
// @name apikey
func main() {
	flag.Parse()
	defer log.Root.Close()

	serror.Service = config.Servicename
	config.File = configFile
	if config.File == "" {
		cfgFile, err := config.GetDefaultConfigfile()
		if err != nil {
			log.Root.Errorf("error getting default config file: %v", err)
			panic("error getting default config file")
		}
		config.File = cfgFile
	}

	log.Root.Infof("using config file: %s", configFile)

	if err := config.Load(); err != nil {
		log.Root.Alertf("can't load config file: %s", err.Error())
		panic("can't load config file")
	}

	serviceConfig = config.Get()
	initConfig()
	initLogging()

	if err := services.InitServices(serviceConfig); err != nil {
		log.Root.Alertf("error creating services: %v", err)
		panic("error creating services")
	}
	log.Root.Info("service is starting")

	var closer io.Closer
	tracer, closer = initJaeger(config.Servicename, serviceConfig.OpenTracing)
	defer closer.Close()

	log.Root.Infof("ssl: %t", serviceConfig.Services.HTTP.Sslport > 0)
	log.Root.Infof("serviceURL: %s", serviceConfig.Services.HTTP.ServiceURL)
	log.Root.Infof("apikey: %s", apiv1.APIKey)
	router, err := apiv1.APIRoutes(serviceConfig, tracer)
	if err != nil {
		errstr := fmt.Sprintf("could not create api routes. %s", err.Error())
		log.Root.Alertf(errstr)
		panic(errstr)
	}

	healthRouter := apiv1.HealthRoutes(serviceConfig, tracer)

	sh := do.MustInvokeNamed[shttp.SHttp](nil, shttp.DoSHTTP)
	sh.StartServers(router, healthRouter)

	log.Root.Info("waiting for clients")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	sh.ShutdownServers()
	log.Root.Info("finished")

	os.Exit(0)
}

// initLogging initialize the logging, especially the gelf logger
func initLogging() {
	var err error
	serviceConfig.Logging.Filename, err = config.ReplaceConfigdir(serviceConfig.Logging.Filename)
	if err != nil {
		log.Root.Errorf("error on config dir: %v", err)
	}
	log.Init(serviceConfig.Logging)
}

// initConfig override the configuration from the service.yaml with the given commandline parameters
func initConfig() {
	if port > 0 {
		serviceConfig.Services.HTTP.Port = port
	}
	if sslport > 0 {
		serviceConfig.Services.HTTP.Sslport = sslport
	}
	if serviceURL != "" {
		serviceConfig.Services.HTTP.ServiceURL = serviceURL
	}
	serviceConfig.Provide()
}

// initJaeger initialize the jaeger (opentracing) component
func initJaeger(servicename string, cnfg config.OpenTracing) (opentracing.Tracer, io.Closer) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: cnfg.Host,
			CollectorEndpoint:  cnfg.Endpoint,
		},
	}
	if (cnfg.Endpoint == "") && (cnfg.Host == "") {
		cfg.Disabled = true
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}
