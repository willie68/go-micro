package bootstrap

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/adapter/outbound/address"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/infrastructure/health"
	"github.com/willie68/go-micro/internal/infrastructure/shttp"
)

func TestInitServices(t *testing.T) {
	inj := do.New()
	cfg := config.Config{
		AddressStorage: address.Config{Type: "internal"},
		HealthSystem:   health.Config{Period: 30, StartDelay: 0},
		HTTP:           shttp.Config{Port: 0, Servicename: "test"},
	}
	assert.NoError(t, InitServices(inj, cfg))
	ShutdownServices(inj)
}
