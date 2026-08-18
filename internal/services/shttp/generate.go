package shttp

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// generateCertificate model
type generateCertificate struct {
	ServiceName  string
	CA           string
	Organization string
	Host         string
	DNSnames     []string
	IPs          []string
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

// GenerateTLSConfig generates the config with self signed certificates
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
		logger.Error(fmt.Sprintf("Unrecognized elliptic curve: %q", gc.EcdsaCurve))
		return nil, err
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to generate private key: %v", err))
		return nil, err
	}

	var notBefore time.Time
	if len(gc.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", gc.ValidFrom)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse creation date: %v", err))
			return nil, err
		}
	}

	template, err := gc.createTemplate(notBefore)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to create certificate template: %v", err))
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, gc.publicKey(priv), priv)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to create certificate: %v", err))
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		logger.Warn(fmt.Sprintf("Unable to marshal private key: %v", err))
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to combine tls key pair: %v", err))
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}

func (gc *generateCertificate) createTemplate(notBefore time.Time) (*x509.Certificate, error) {
	notAfter := notBefore.Add(gc.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to generate serial number: %v", err))
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{gc.Organization},
			CommonName:   gc.ServiceName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, sip := range gc.IPs {
		if ip := net.ParseIP(sip); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		}
	}

	template.DNSNames = append(template.DNSNames, gc.DNSnames...)

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
	return &template, nil
}

func (s *SHttp) TLSFromFiles() (*tls.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair(s.cfn.Certificate, s.cfn.Key)
	if err != nil {
		return nil, errors.Wrap(err, "error creating X509 Pair")
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}
