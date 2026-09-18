package shttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestProvideAndLifecycle(t *testing.T) {
	inj := do.New()
	do.ProvideValue(inj, Config{
		Servicename: "coverage",
		Port:        0,
		ServiceURL:  "http://127.0.0.1",
	})
	assert.NoError(t, Provide(inj))
	sh := do.MustInvoke[*SHttp](inj)
	assert.False(t, sh.Started())

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	sh.StartServers(router, router)
	time.Sleep(50 * time.Millisecond)
	assert.True(t, sh.Started())
	sh.ShutdownServers()
	assert.False(t, sh.Started())
}

func TestStartHTTPSGeneratedCert(t *testing.T) {
	sh, err := NewSHttp(Config{
		Servicename: "coverage",
		Port:        0,
		Sslport:     0,
		ServiceURL:  "https://localhost",
	})
	assert.NoError(t, err)
	assert.False(t, sh.useSSL)

	sh, err = NewSHttp(Config{
		Servicename: "coverage",
		Port:        0,
		Sslport:     18443,
		ServiceURL:  "https://localhost",
	})
	assert.NoError(t, err)
	assert.True(t, sh.useSSL)
	router := chi.NewRouter()
	sh.StartServers(router, router)
	time.Sleep(80 * time.Millisecond)
	assert.True(t, sh.Started())
	sh.ShutdownServers()
}
