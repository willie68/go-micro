package services

import (
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/sconfig"
	"github.com/willie68/go-micro/internal/services/shttp"
)

var(
	logger = logging.New().WithName("services")
	healthService *health.SHealth
)

// InitServices initialise the service system
func InitServices(cfg config.Config) error {
	err := InitHelperServices(cfg)
	if err != nil {
		return err
	}

	// here you can add more services
	s, err := sconfig.NewSConfig()
	if err != nil {
		return err
	}

	healthService.Register(s)

	return InitRESTService(cfg)
}

// InitHelperServices initialise the helper services like Healthsystem
func InitHelperServices(cfg config.Config) error {
	var err error
	healthService, err = health.NewHealthSystem(cfg.Service.HealthSystem)
	if err != nil {
		return err
	}
	return nil
}

// InitRESTService initialise REST Services
func InitRESTService(cfg config.Config) error {
	_, err := shttp.NewSHttp(cfg.Service.HTTP, cfg.Service.CA)
	if err != nil {
		return err
	}
	return nil
}
