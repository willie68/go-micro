// Package main this is the entry point into the service
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/samber/do/v2"
	_ "github.com/willie68/go-micro/docs"
	"github.com/willie68/go-micro/internal"
	"github.com/willie68/go-micro/internal/apiv1"
	"github.com/willie68/go-micro/internal/serror"
	"github.com/willie68/go-micro/internal/services/shttp"

	config "github.com/willie68/go-micro/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	log "github.com/willie68/go-micro/internal/services/logging"

	flag "github.com/spf13/pflag"
)

var (
	configFile    string
	serviceConfig config.Config
	c             chan os.Signal
)

func init() {
	// variables for parameter override
	log.Root.Info("init service")
	flag.StringVarP(&configFile, "config", "c", config.File, "this is the path and filename to the config file")
}

// @title			GoMicro service API
// @version		1.0
// @description	The GoMicro service is a template for microservices written in go.
// @BasePath		/api/v1
// @in				header
func main() {
	inj := do.New()
	flag.Parse()

	serror.Service = config.Servicename
	config.File = configFile
	if config.File == "" {
		cfgFile, err := config.GetDefaultConfigfile()
		if err != nil {
			panic(fmt.Sprintf("error getting default config file: %v", err))
		}
		config.File = cfgFile
	}

	if err := config.Load(); err != nil {
		log.Root.Warn(fmt.Sprintf("can't load config file: %v", err))
		panic("can't load config file")
	}

	log.Init(config.Get().Logging, config.Servicename)
	log.Root.Info(fmt.Sprintf("using config file: %s: '%s'", configFile, config.YAML()))

	serviceConfig = config.Get()
	serviceConfig.Provide(inj)

	if err := internal.InitServices(inj, serviceConfig); err != nil {
		log.Root.Warn(fmt.Sprintf("error creating services: %v", err))
		panic("error creating services")
	}
	log.Root.Info("service is starting")

	tp, err := initOpenTelemetry(config.Servicename, serviceConfig.OpenTelemetry)
	if err != nil {
		panic(fmt.Sprintf("opentelemetry init failed: %v", err))
	}

	log.Root.Info(fmt.Sprintf("ssl: %t", serviceConfig.HTTP.Sslport > 0))
	log.Root.Info(fmt.Sprintf("serviceURL: %s", serviceConfig.HTTP.ServiceURL))
	router, err := apiv1.APIRoutes(inj, serviceConfig)
	if err != nil {
		errstr := fmt.Sprintf("could not create api routes. %s", err.Error())
		log.Root.Warn(errstr)
		panic(errstr)
	}

	healthRouter := apiv1.HealthRoutes(inj, serviceConfig)

	sh := do.MustInvoke[shttp.SHttp](inj)
	sh.StartServers(router, healthRouter)

	log.Root.Info("waiting for clients")
	c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	sh.ShutdownServers()
	tp.Shutdown(context.Background())

	log.Root.Info("finished")
	os.Exit(0)
}

// initLogging initialize the logging, especially the gelf logger
func initLogging() {
	var err error
	serviceConfig.Logging.Filename, err = config.ReplaceConfigdir(serviceConfig.Logging.Filename)
	if err != nil {
		log.Root.Error(fmt.Sprintf("error on config dir: %v", err))
	}
	log.Init(serviceConfig.Logging, "gomicro")
}

// initOpenTelemetry initialize the opentelemetry component
func initOpenTelemetry(servicename string, cnfg config.OpenTelemetry) (*sdktrace.TracerProvider, error) {
	endpoint := strings.TrimSpace(cnfg.Endpoint)
	if endpoint == "" {
		return nil, nil
	}

	clientOptions := []otlptracehttp.Option{}
	if strings.Contains(endpoint, "://") {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		clientOptions = append(clientOptions, otlptracehttp.WithEndpoint(u.Host))
		if u.Path != "" && u.Path != "/" {
			clientOptions = append(clientOptions, otlptracehttp.WithURLPath(strings.TrimPrefix(u.Path, "/")))
		}
		if strings.EqualFold(u.Scheme, "http") {
			clientOptions = append(clientOptions, otlptracehttp.WithInsecure())
		}
	} else {
		clientOptions = append(clientOptions, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(context.Background(), clientOptions...)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(servicename),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}
