package shttp

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/caservice"
)

var logger = logging.New().WithName("shttp")

// SHttp a service encapsulating http and https server
type SHttp struct {
	cfn     Config
	cfa     caservice.Config
	useSSL  bool
	sslsrv  *http.Server
	srv     *http.Server
	Started bool
}

// NewSHttp creates a new shttp service
func NewSHttp(cfn Config, cfgCa caservice.Config) (*SHttp, error) {
	sh := SHttp{
		cfn:     cfn,
		cfa:     cfgCa,
		Started: false,
	}
	sh.init()

	do.ProvideValue[SHttp](nil, sh)

	return &sh, nil
}

func (s *SHttp) init() {
	if s.cfn.Sslport > 0 {
		s.useSSL = true
	}
	s.Started = false
}

// StartServers starting all needed http servers
func (s *SHttp) StartServers(router, healthRouter *chi.Mux) {
	if s.useSSL {
		s.startHTTPSServer(router)
		s.startHTTPServer(healthRouter)
	} else {
		s.startHTTPServer(router)
	}
	s.Started = true
}

// ShutdownServers shutting all servers down
func (s *SHttp) ShutdownServers() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Errorf("shutdown http server error: %v", err)
	}
	if s.useSSL {
		if err := s.sslsrv.Shutdown(ctx); err != nil {
			logger.Errorf("shutdown https server error: %v", err)
		}
	}
	s.Started = false
}

func (s *SHttp) startHTTPSServer(router *chi.Mux) {
	var tlsConfig *tls.Config
	var err error
	if s.cfa.UseCA {
		tlsConfig, err = s.GetTLSConfig()
		if err != nil {
			logger.Alertf("could not create tls config. %s", err.Error())
			panic(-1)
		}
	} else {
		if s.cfn.Certificate != "" && s.cfn.Key != "" {
			// using the files provided by config
			tlsConfig, err = s.TLSFromFiles()
			if err != nil {
				logger.Alertf("could not create tls config. %s", err.Error())
				panic(-1)
			}
		} else {
			// generating our own certificate
			h := s.cfn.ServiceURL
			ul, err := url.Parse(h)
			if err == nil {
				h = ul.Hostname()
			}
			gc := generateCertificate{
				ServiceName: s.cfa.Servicename,
				CA:          s.cfa.URL,
				Host:        h,
				ValidFor:    10 * 365 * 24 * time.Hour,
				IsCA:        false,
				EcdsaCurve:  "P384",
				Ed25519Key:  false,
				DNSnames:    s.cfn.DNSNames,
				IPs:         s.cfn.IPAddresses,
			}
			tlsConfig, err = gc.GenerateTLSConfig()
			if err != nil {
				logger.Alertf("could not create tls config. %s", err.Error())
				panic(-1)
			}
		}
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
		logger.Infof("starting https server on address: %s", s.sslsrv.Addr)
		if err := s.sslsrv.ListenAndServeTLS("", ""); err != nil {
			logger.Alertf("error starting server: %s", err.Error())
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
		logger.Infof("starting http server on address: %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil {
			logger.Alertf("error starting server: %s", err.Error())
		}
	}()
}
