package shttp

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/config"
	log "github.com/willie68/go-micro/internal/logging"
	mv "github.com/willie68/micro-vault/pkg/client"
)

const (
	// DoSHTTP naming constant for dependency injection
	DoSHTTP = "shttp"
)

// SHttp a service encapsulating http and https server
type SHttp struct {
	cfn     config.Config
	useSSL  bool
	sslsrv  *http.Server
	srv     *http.Server
	Started bool
}

// NewSHttp creates a new shttp service
func NewSHttp(cfn config.Config) (*SHttp, error) {
	sh := SHttp{
		cfn:     cfn,
		Started: false,
	}
	sh.init()

	do.ProvideNamedValue[SHttp](nil, DoSHTTP, sh)

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
		log.Logger.Errorf("shutdown http server error: %v", err)
	}
	if s.useSSL {
		if err := s.sslsrv.Shutdown(ctx); err != nil {
			log.Logger.Errorf("shutdown https server error: %v", err)
		}
	}
	s.Started = false
}

func (s *SHttp) startHTTPSServer(router *chi.Mux) {
	var tlsConfig *tls.Config
	var err error
	if s.cfn.CA.UseCA {
		tlsConfig, err = s.GetTLSConfig()
		if err != nil {
			log.Logger.Alertf("could not create tls config. %s", err.Error())
		}
	} else {
		h := s.cfn.ServiceURL
		ul, err := url.Parse(h)
		if err == nil {
			h = ul.Hostname()
		}
		gc := generateCertificate{
			ServiceName: config.Servicename,
			CA:          s.cfn.CA.URL,
			Host:        h,
			ValidFor:    10 * 365 * 24 * time.Hour,
			IsCA:        false,
			EcdsaCurve:  "P384",
			Ed25519Key:  false,
		}
		tlsConfig, err = gc.GenerateTLSConfig()
		if err != nil {
			log.Logger.Alertf("could not create tls config. %s", err.Error())
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

// generateCertificate model
type generateCertificate struct {
	ServiceName  string
	CA           string
	Organization string
	Host         string
	ValidFrom    string
	ValidFor     time.Duration
	IsCA         bool
	RSABits      int
	EcdsaCurve   string
	Ed25519Key   bool
}

func (gc *generateCertificate) publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

// GenerateTLSConfig generates the config
func (gc *generateCertificate) GenerateTLSConfig() (*tls.Config, error) {
	var priv any
	var err error
	switch gc.EcdsaCurve {
	case "":
		if gc.Ed25519Key {
			_, priv, err = ed25519.GenerateKey(rand.Reader)
		} else {
			priv, err = rsa.GenerateKey(rand.Reader, gc.RSABits)
		}
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		log.Logger.Fatalf("Unrecognized elliptic curve: %q", gc.EcdsaCurve)
		return nil, err
	}
	if err != nil {
		log.Logger.Fatalf("Failed to generate private key: %v", err)
		return nil, err
	}

	var notBefore time.Time
	if len(gc.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", gc.ValidFrom)
		if err != nil {
			log.Logger.Fatalf("Failed to parse creation date: %v", err)
			return nil, err
		}
	}

	notAfter := notBefore.Add(gc.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Logger.Fatalf("Failed to generate serial number: %v", err)
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{gc.Organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(gc.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if gc.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, gc.publicKey(priv), priv)
	if err != nil {
		log.Logger.Fatalf("Failed to create certificate: %v", err)
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Logger.Fatalf("Unable to marshal private key: %v", err)
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Logger.Fatalf("Failed to combine tls key pair: %v", err)
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}

// GetTLSConfig generates the tls config, getting certificate from ca service
func (s *SHttp) GetTLSConfig() (*tls.Config, error) {
	var err error
	subj := pkix.Name{
		CommonName: config.Servicename,
	}
	rawSubj := subj.ToRDNSequence()

	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		return nil, err
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	h := s.cfn.ServiceURL
	ul, err := url.Parse(h)
	if err == nil {
		h = ul.Hostname()
	}
	if ip := net.ParseIP(h); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, h)
	}

	cli, err := mv.LoginClient(s.cfn.CA.AccessKey, s.cfn.CA.Secret, s.cfn.CA.URL)
	if err != nil {
		return nil, err
	}

	crt, err := cli.Certificate(template)
	if err != nil {
		return nil, err
	}

	prv, err := cli.PrivateKey()
	if err != nil {
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(prv)
	if err != nil {
		log.Logger.Fatalf("Unable to marshal private key: %v", err)
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}
