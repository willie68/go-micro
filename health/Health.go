package health

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	clog "wkla.no-ip.biz/go-micro/logging"
)

var myhealthy bool

/*
This is the healtchcheck you will have to provide.
*/
func check(tracer opentracing.Tracer) (bool, string) {
	// TODO implement here your healthcheck.
	myhealthy = !myhealthy
	message := ""
	if myhealthy {
		clog.Logger.Info("healthy")
	} else {
		clog.Logger.Info("not healthy")
		message = "ungesund"
	}
	return myhealthy, message
}

//##### template internal functions for processing the healthchecks #####
var healthmessage string
var healthy bool
var lastChecked time.Time
var period int

// CheckConfig configuration for the healthcheck system
type CheckConfig struct {
	Period int
}

// Msg a health message
type Msg struct {
	Message   string `json:"message"`
	LastCheck string `json:"lastCheck,omitempty"`
}

// InitHealthSystem initialise the complete health system
func InitHealthSystem(config CheckConfig, tracer opentracing.Tracer) {
	period = config.Period
	clog.Logger.Infof("healthcheck starting with period: %d seconds", period)
	healthmessage = "service starting"
	healthy = false
	doCheck(tracer)
	go func() {
		background := time.NewTicker(time.Second * time.Duration(period))
		for _ = range background.C {
			doCheck(tracer)
		}
	}()
}

/*
internal function to process the health check
*/
func doCheck(tracer opentracing.Tracer) {
	var msg string
	healthy, msg = check(tracer)
	if !healthy {
		healthmessage = msg
	} else {
		healthmessage = ""
	}
	lastChecked = time.Now()
}

/*
Routes getting all routes for the health endpoint
*/
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthz", GetHealthyEndpoint)
	router.Get("/readyz", GetReadinessEndpoint)
	return router
}

/*
GetHealthyEndpoint liveness probe
*/
func GetHealthyEndpoint(response http.ResponseWriter, req *http.Request) {
	render.Status(req, http.StatusOK)
	render.JSON(response, req, Msg{
		Message: fmt.Sprintf("service started"),
	})
}

/*
GetReadinessEndpoint is this service ready for taking requests, e.g. formaly known as health checks
*/
func GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	t := time.Now()
	if t.Sub(lastChecked) > (time.Second * time.Duration(2*period)) {
		healthy = false
		healthmessage = "Healthcheck not running"
	}
	if healthy {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Msg{
			Message:   "service up and running",
			LastCheck: lastChecked.String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, Msg{
			Message:   fmt.Sprintf("service is unavailable: %s", healthmessage),
			LastCheck: lastChecked.String(),
		})
	}
}

/*
sendMessage sending a span message to tracer
*/
func sendMessage(tracer opentracing.Tracer, message string) {
	span := tracer.StartSpan("say-hello")
	println(message)
	span.Finish()
}
