package shttp

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/crypt"
	log "github.com/willie68/go-micro/internal/logging"
)

const (
	// DoSHTTP naming constant for dependency injection
	DoSHTTP = "shttp"
)

type SHttp struct {
	cfn    config.Config
	useSSL bool
	sslsrv *http.Server
	srv    *http.Server
}

func NewSHttp(cfn config.Config) (*SHttp, error) {
	sh := SHttp{
		cfn: cfn,
	}
	sh.init()

	do.ProvideNamedValue[SHttp](nil, DoSHTTP, sh)

	return &sh, nil
}

func (s *SHttp) init() {
	if s.cfn.Sslport > 0 {
		s.useSSL = true
	}
}

// StartServers starting all needed http servers
func (s *SHttp) StartServers(router, healthRouter *chi.Mux) {
	if s.useSSL {
		s.startHTTPSServer(router)
		s.startHTTPServer(healthRouter)
	} else {
		s.startHTTPServer(router)
	}
}

// ShutdownServers shutting all servers down
func (s *SHttp) ShutdownServers() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Logger.Errorf("shutdown http server error: %v", err)
	}
	if s.useSSL {
		if err := s.sslsrv.Shutdown(ctx); err != nil {
			log.Logger.Errorf("shutdown https server error: %v", err)
		}
	}
}

func (s *SHttp) startHTTPSServer(router *chi.Mux) {
	gc := crypt.GenerateCertificate{
		Organization: "MCS",
		Host:         "127.0.0.1",
		ValidFor:     10 * 365 * 24 * time.Hour,
		IsCA:         false,
		EcdsaCurve:   "P384",
		Ed25519Key:   false,
	}
	tlsConfig, err := gc.GenerateTLSConfig()
	if err != nil {
		log.Logger.Alertf("could not create tls config. %s", err.Error())
	}
	s.sslsrv = &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(s.cfn.Sslport),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
		TLSConfig:    tlsConfig,
	}
	go func() {
		log.Logger.Infof("starting https server on address: %s", s.sslsrv.Addr)
		if err := s.sslsrv.ListenAndServeTLS("", ""); err != nil {
			log.Logger.Alertf("error starting server: %s", err.Error())
		}
	}()
}

func (s *SHttp) startHTTPServer(router *chi.Mux) {
	// own http server for the healthchecks
	s.srv = &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(s.cfn.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	go func() {
		log.Logger.Infof("starting http server on address: %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil {
			log.Logger.Alertf("error starting server: %s", err.Error())
		}
	}()
}
