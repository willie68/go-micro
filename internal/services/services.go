package services

import (
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/adrsvc/adrfact"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/shttp"
)

var (
	logger = logging.New().WithName("services")
)

// InitServices initialise the service system
func InitServices(cfg config.Config) error {
	logger.Debug("initialise services")
	err := InitHelperServices(cfg)
	if err != nil {
		return err
	}

	// here you can add more services
	s, err := adrfact.New(cfg.AddressStorage)
	if err != nil {
		return err
	}

	err = health.Register(s)
	if err != nil {
		return err
	}

	return InitRESTService(cfg)
}

// InitHelperServices initialise the helper services like Healthsystem
func InitHelperServices(cfg config.Config) error {
	var err error
	_, err = health.NewHealthSystem(cfg.HealthSystem)
	return err
}

// InitRESTService initialise REST Services
func InitRESTService(cfg config.Config) error {
	_, err := shttp.NewSHttp(cfg.HTTP, cfg.CA)
	return err
}
