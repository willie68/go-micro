package client

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/adapter/inbound/http/apiv1"
	"github.com/willie68/go-micro/internal/bootstrap"
	"github.com/willie68/go-micro/internal/config"
)

var (
	srvStarted bool
	sh         SHttp
	cfg        config.Config
)

type SHttp interface {
	StartServers(router, healthRouter *chi.Mux)
	ShutdownServers()
	Started() bool
}

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
		if err := bootstrap.InitServices(inj, cfg); err != nil {
			panic("error creating services")
		}

		sh = do.MustInvokeAs[SHttp](inj)
	}
	if !sh.Started() {
		router, err := apiv1.APIRoutes(inj, cfg)
		if err != nil {
			errstr := fmt.Sprintf("could not create api routes. %s", err.Error())
			panic(errstr)
		}

		healthRouter := apiv1.HealthRoutes(inj, cfg)
		sh.StartServers(router, healthRouter)

		time.Sleep(1 * time.Second)
	}
}

func TestStartServer(t *testing.T) {
	inj := do.New()
	ast := assert.New(t)

	StartServer(inj)

	ast.NotNil(sh)
	ast.True(sh.Started())
}
