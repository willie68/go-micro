package services

import (
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/grpc"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/sconfig"
	"github.com/willie68/go-micro/internal/services/shttp"
)

var (
	logger        = logging.New().WithName("services")
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

	InitGRPCService(cfg)
	return InitRESTService(cfg)
}

// InitHelperServices initialise the helper services like Healthsystem
func InitHelperServices(cfg config.Config) error {
	var err error
	healthService, err = health.NewHealthSystem(cfg.Services.HealthSystem)
	if err != nil {
		return err
	}
	return nil
}

func InitGRPCService(cfg config.Config) {
	grpc.NewGRPC(cfg.Services.GRPC)
}

// InitRESTService initialise REST Services
func InitRESTService(cfg config.Config) error {
	_, err := shttp.NewSHttp(cfg.Services.HTTP, cfg.Services.CA)
	if err != nil {
		return err
	}
	return nil
}
