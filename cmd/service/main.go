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
	configFile    string
	serviceConfig config.Config
	tracer        opentracing.Tracer
	c             chan os.Signal
)

func init() {
	// variables for parameter override
	log.Root.Info("init service")
	flag.StringVarP(&configFile, "config", "c", config.File, "this is the path and filename to the config file")
}

// @title GoMicro service API
// @version 1.0
// @description The GoMicro service is a template for microservices written in go.
// @BasePath /api/v1
// @in header
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
	serviceConfig.Provide()
	initLogging()

	if err := services.InitServices(serviceConfig); err != nil {
		log.Root.Alertf("error creating services: %v", err)
		panic("error creating services")
	}
	log.Root.Info("service is starting")

	var closer io.Closer
	tracer, closer = initJaeger(config.Servicename, serviceConfig.OpenTracing)
	defer closer.Close()

	log.Root.Infof("ssl: %t", serviceConfig.HTTP.Sslport > 0)
	log.Root.Infof("serviceURL: %s", serviceConfig.HTTP.ServiceURL)
	router, err := apiv1.APIRoutes(serviceConfig, tracer)
	if err != nil {
		errstr := fmt.Sprintf("could not create api routes. %s", err.Error())
		log.Root.Alertf(errstr)
		panic(errstr)
	}

	healthRouter := apiv1.HealthRoutes(serviceConfig, tracer)

	sh := do.MustInvoke[shttp.SHttp](nil)
	sh.StartServers(router, healthRouter)

	log.Root.Info("waiting for clients")
	c = make(chan os.Signal, 1)
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
