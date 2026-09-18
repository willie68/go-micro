package bootstrap

import (
	"github.com/samber/do/v2"
	adrsvc "github.com/willie68/go-micro/internal/adapter/outbound/address"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/domain/addresses"
	"github.com/willie68/go-micro/internal/infrastructure/health"
	"github.com/willie68/go-micro/internal/infrastructure/logging"
	"github.com/willie68/go-micro/internal/infrastructure/shttp"
)

var (
	logger = logging.New("services")
)

// InitServices initialise the service system
func InitServices(inj do.Injector, cfg config.Config) error {
	logger.Debug("initialise services")
	do.ProvideValue(inj, cfg)
	do.ProvideValue(inj, cfg.AddressStorage)
	do.ProvideValue(inj, cfg.HealthSystem)
	do.ProvideValue(inj, cfg.HTTP)

	err := InitHelperServices(inj)
	if err != nil {
		return err
	}

	// here you can add more services
	err = adrsvc.Provide(inj)
	if err != nil {
		return err
	}

	addresses.Provide(inj)

	return InitRESTService(inj)
}

// InitHelperServices initialise the helper services like Healthsystem
func InitHelperServices(inj do.Injector) error {
	return health.Provide(inj)
}

// InitRESTService initialise REST Services
func InitRESTService(inj do.Injector) error {
	return shttp.Provide(inj)
}

// ShutdownServices shutting down all services, that support do.Shutdowner interface
func ShutdownServices(inj do.Injector) {
	inj.Shutdown()
}
