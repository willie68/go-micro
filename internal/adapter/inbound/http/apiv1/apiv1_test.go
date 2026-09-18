package apiv1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/infrastructure/health"
	"github.com/willie68/go-micro/pkg/pmodel"
)

type routeServiceName struct{}

func (routeServiceName) ServiceName() string { return "go-micro" }

func TestToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer abc")
	tk, err := token(req)
	assert.NoError(t, err)
	assert.Equal(t, "abc", tk)
}

func TestAPIRoutesAndHealthRoutes(t *testing.T) {
	inj := do.New()
	stg := &fakeAddress{items: map[string]pmodel.Address{}}
	do.ProvideValue(inj, stg)
	hsvc, err := health.NewHealthSystem(inj, health.Config{Period: 30, StartDelay: 0})
	assert.NoError(t, err)
	do.ProvideValue(inj, hsvc)
	do.ProvideValue(inj, routeServiceName{})

	cfg := config.Config{
		Metrics: config.Metrics{Enable: true},
	}
	router, err := APIRoutes(inj, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, router)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/addresses/", nil)
	req.Header.Set("tenant", "t1")
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	hr := HealthRoutes(inj, config.Config{Metrics: config.Metrics{Enable: true}, Profiling: config.Profiling{Enable: true}})
	assert.NotNil(t, hr)
}
