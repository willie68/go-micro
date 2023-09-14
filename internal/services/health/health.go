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
	log "github.com/willie68/go-micro/internal/logging"
)

const DoSHealth = "shealth"

// Config configuration for the healthcheck system
type Config struct {
	// Period in seconds, when all health services should run
	Period int `yaml:"period"`
	// StartDelay an optional starting delay, after starting the service
	StartDelay int `yaml:"startdelay"`
}

type Check interface {
	// CheckName shoulld return the name of this healthcheck. The name should be unique.
	CheckName() string
	// Check proceed a check and return state, true for healthy or false and an optional error, if the healthcheck fails
	Check() (bool, error)
}

// SHealthCheck this is the healthcheck service
type SHealth struct {
	cfg         Config
	healthy     bool
	readyz      bool
	messages    []string
	checks      []Check
	lastChecked time.Time
	reg         sync.Mutex
}

// Msg a health message
type Message struct {
	Messages  []string `json:"messages"`
	LastCheck string   `json:"lastCheck,omitempty"`
}

// NewHealthSystem initialize the complete health system
func NewHealthSystem(config Config) (*SHealth, error) {
	shealth := SHealth{
		cfg:     config,
		healthy: false,
		checks:  make([]Check, 0),
		reg:     sync.Mutex{},
	}
	err := shealth.Init()
	if err != nil {
		return nil, err
	}
	do.ProvideNamedValue[*SHealth](nil, DoSHealth, &shealth)
	return &shealth, nil
}

func (h *SHealth) Init() error {
	log.Logger.Infof("healthcheck starting with period: %d seconds", h.cfg.Period)
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
func (h *SHealth) checkHealthCheckTimer() {
	t := time.Now()
	if t.Sub(h.lastChecked) > (time.Second * time.Duration(2*h.cfg.Period)) {
		h.readyz = false
		h.messages = []string{"health check not running"}
		if t.Sub(h.lastChecked) > (time.Second * time.Duration(4*h.cfg.Period)) {
			log.Logger.Error("panic: health check is not running anymore")
			panic("panic: health check is not running anymore")
		}
	}
}

func (h *SHealth) Register(check Check) {
	h.reg.Lock()
	for x, c := range h.checks {
		if c.CheckName() == check.CheckName() {
			h.checks[x] = check
			return
		}
	}
	h.checks = append(h.checks, check)
	defer h.reg.Unlock()
}

func (h *SHealth) Message() Message {
	return Message{
		LastCheck: h.lastChecked.String(),
		Messages:  h.messages,
	}
}

// doCheck internal function to process the health check
func (h *SHealth) doCheck() {
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

// Handler is the default handler factory for HTTP requests against the healthsystem
type Handler struct {
	health *SHealth
}

func NewHealthHandler() api.Handler {
	return &Handler{
		health: do.MustInvokeNamed[*SHealth](nil, DoSHealth),
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
	if h.health.healthy {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.JSON(response, req, h.health.Message())
}

// HeadLivenessEndpoint liveness probe
func (h *Handler) HeadLivenessEndpoint(response http.ResponseWriter, req *http.Request) {
	if h.health.healthy {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.NoContent(response, req)
}

// GetReadinessEndpoint is this service ready for taking requests, e.g. formerly known as health checksfunc GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
func (h *Handler) GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	h.health.checkHealthCheckTimer()
	if h.health.readyz {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Message{
			Messages:  []string{"main: service up and running"},
			LastCheck: h.health.lastChecked.String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, h.health.Message())
	}
}

// HeadReadinessEndpoint is this service ready for taking requests, e.g. formaly known as health checks
func (h *Handler) HeadReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	h.health.checkHealthCheckTimer()
	if h.health.readyz {
		render.Status(req, http.StatusOK)
	} else {
		render.Status(req, http.StatusServiceUnavailable)
	}
	render.NoContent(response, req)
}
