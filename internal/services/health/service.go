package health

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/api"
	"github.com/willie68/go-micro/internal/logging"
)

//go:generate mockery --name=SHealth --outpkg=servicemocks --filename=service.go --with-expecter
type SHealth interface {
	Healthy() bool
	Readyz() bool
	LastChecked() time.Time
	Message() Message
	Register(check Check)
	Unregister(checkname string) bool

	CheckHealthCheckTimer()
}

var logger = logging.New().WithName("health")

// Service  this is the healthcheck service
type Service struct {
	cfg         Config
	healthy     bool
	readyz      bool
	messages    []string
	checks      []Check
	lastChecked time.Time
	reg         sync.Mutex
}

var _ SHealth = &Service{}

// Message a health message
type Message struct {
	Messages  []string `json:"messages"`
	LastCheck string   `json:"lastCheck,omitempty"`
}

// NewHealthSystem initialize the complete health system
func NewHealthSystem(config Config) (SHealth, error) {
	shealth := Service{
		cfg:     config,
		healthy: false,
		checks:  make([]Check, 0),
		reg:     sync.Mutex{},
	}
	err := shealth.Init()
	if err != nil {
		return nil, err
	}
	do.ProvideValue[SHealth](nil, &shealth)
	return &shealth, nil
}

// Init initialise the health system
func (h *Service) Init() error {
	logger.Infof("healthcheck starting with period: %d seconds", h.cfg.Period)
	h.messages = make([]string, 0)
	h.messages = append(h.messages, "service starting")
	h.readyz = false
	h.doCheck()
	h.lastChecked = time.Now()
	go func() {
		if h.cfg.StartDelay > 0 {
			time.Sleep(time.Duration(h.cfg.StartDelay) * time.Second)
		}
		go func() {
			if h.cfg.Period > 0 {
				background := time.NewTicker(time.Second * time.Duration(h.cfg.Period))
				for range background.C {
					h.doCheck()
				}
			}
		}()
	}()
	return nil
}

// checking if the health system (namly the timer task) is working or stopped
func (h *Service) CheckHealthCheckTimer() {
	t := time.Now()
	if t.Sub(h.lastChecked) > (time.Second * time.Duration(2*h.cfg.Period)) {
		h.readyz = false
		h.messages = []string{"health check not running"}
		if t.Sub(h.lastChecked) > (time.Second * time.Duration(4*h.cfg.Period)) {
			logger.Error("panic: health check is not running anymore")
			panic("panic: health check is not running anymore")
		}
	}
}

// Register register a new healthcheck. If a healthcheck with the same name is already present, this will be overwritten
// Otherwise the new healthcheck will be appended
func (h *Service) Register(check Check) {
	h.reg.Lock()
	defer h.reg.Unlock()
	for x, c := range h.checks {
		if c.CheckName() == check.CheckName() {
			h.checks[x] = check
			return
		}
	}
	h.checks = append(h.checks, check)
}

// Unregister unregister a healthcheck. Return true if the healthcheck can be unregistered otherwise false
func (h *Service) Unregister(checkname string) bool {
	h.reg.Lock()
	defer h.reg.Unlock()
	for x := len(h.checks) - 1; x >= 0; x-- {
		c := h.checks[x]
		if c.CheckName() == checkname {
			h.checks = append(h.checks[:x], h.checks[x+1:]...)
			return true
		}
	}
	return false
}

// Message return a health message from the last healthcheck
func (h *Service) Message() Message {
	return Message{
		LastCheck: h.lastChecked.String(),
		Messages:  h.messages,
	}
}

// doCheck internal function to process the health check
func (h *Service) doCheck() {
	h.lastChecked = time.Now()
	h.messages = make([]string, 0)
	healthy := true
	h.reg.Lock()
	for _, c := range h.checks {
		ok, err := c.Check()
		if !ok {
			healthy = false
			h.messages = append(h.messages, fmt.Sprintf("%s: %s", c.CheckName(), err.Error()))
		}
	}
	defer h.reg.Unlock()
	h.healthy = healthy
	if healthy {
		h.readyz = true
	}
}

// Healthy return the actual health state
func (h *Service) Healthy() bool {
	return h.healthy
}

// Healthy return the actual health state
func (h *Service) Readyz() bool {
	return h.readyz
}

// Healthy return the actual health state
func (h *Service) LastChecked() time.Time {
	return h.lastChecked
}

// Handler is the default handler factory for HTTP requests against the healthsystem
type Handler struct {
	health SHealth
}

// NewHealthHandler creates a new healthhandler for a REST interface
func NewHealthHandler() api.Handler {
	return &Handler{
		health: do.MustInvoke[SHealth](nil),
	}
}

// Routes getting all routes for the health endpoint
func (h *Handler) Routes() (string, *chi.Mux) {
	router := chi.NewRouter()
	router.Get("/livez", h.GetLivenessEndpoint)
	router.Get("/readyz", h.GetReadinessEndpoint)
	router.Head("/livez", h.HeadLivenessEndpoint)
	router.Head("/readyz", h.HeadReadinessEndpoint)
	return "/", router
}

// GetLivenessEndpoint liveness probe
func (h *Handler) GetLivenessEndpoint(response http.ResponseWriter, req *http.Request) {
	if h.health.Healthy() {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.JSON(response, req, h.health.Message())
}

// HeadLivenessEndpoint liveness probe
func (h *Handler) HeadLivenessEndpoint(response http.ResponseWriter, req *http.Request) {
	if h.health.Healthy() {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.NoContent(response, req)
}

// GetReadinessEndpoint is this service ready for taking requests, e.g. formerly known as health checksfunc GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
func (h *Handler) GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	h.health.CheckHealthCheckTimer()
	if h.health.Readyz() {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Message{
			Messages:  []string{"main: service up and running"},
			LastCheck: h.health.LastChecked().String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, h.health.Message())
	}
}

// HeadReadinessEndpoint is this service ready for taking requests, e.g. formaly known as health checks
func (h *Handler) HeadReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	h.health.CheckHealthCheckTimer()
	if h.health.Readyz() {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.NoContent(response, req)
}
