package client

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/apiv1"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/services"
	"github.com/willie68/go-micro/internal/services/shttp"
)

var (
	srvStarted bool
	sh         *shttp.SHttp
	cfg        config.Config
)

func StartServer(inj do.Injector) {
	if sh == nil {
		fmt.Println("starting server")
		_ = os.Chdir("../../")
		// loading the config file
		config.File = "./testdata/service_local.yaml"
		err := config.Load()
		if err != nil {
			panic("can't load local config")
		}

		cfg = config.Get()
		cfg.Provide(inj)
		if err := services.InitServices(inj, cfg); err != nil {
			panic("error creating services")
		}

		s := do.MustInvoke[shttp.SHttp](inj)
		sh = &s
	}
	if !sh.Started {
		router, err := apiv1.APIRoutes(inj, cfg, nil)
		if err != nil {
			errstr := fmt.Sprintf("could not create api routes. %s", err.Error())
			panic(errstr)
		}

		healthRouter := apiv1.HealthRoutes(inj, cfg, nil)
		sh.StartServers(router, healthRouter)

		time.Sleep(1 * time.Second)
	}
}

func TestStartServer(t *testing.T) {
	inj := do.New()
	ast := assert.New(t)

	StartServer(inj)

	ast.NotNil(sh)
	ast.True(sh.Started)
}
