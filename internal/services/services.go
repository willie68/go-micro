package services

import (
	"github.com/samber/do/v2"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/shttp"
)

var (
	logger = logging.New("services")
)

// Service is the standard service interface
type Service interface {
	Init() error
	Shutdown() error
}

// InitServices initialise the service system
func InitServices(inj do.Injector, cfg config.Config) error {
	logger.Debug("initialise services")
	err := InitHelperServices(inj, cfg)
	if err != nil {
		return err
	}

	// here you can add more services
	err = adrsvc.New(inj, cfg.AddressStorage)
	if err != nil {
		return err
	}

	return InitRESTService(inj, cfg)
}

// InitHelperServices initialise the helper services like Healthsystem
func InitHelperServices(inj do.Injector, cfg config.Config) error {
	var err error
	_, err = health.NewHealthSystem(inj, cfg.HealthSystem)
	return err
}

// InitRESTService initialise REST Services
func InitRESTService(inj do.Injector, cfg config.Config) error {
	_, err := shttp.NewSHttp(inj, cfg.HTTP, cfg.CA)
	return err
}

func ShutdownServices(inj do.Injector) {
	inj.Shutdown()
}
