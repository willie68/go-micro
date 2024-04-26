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
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	mv "github.com/willie68/micro-vault/pkg/client"
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
		logger.Fatalf("Unrecognized elliptic curve: %q", gc.EcdsaCurve)
		return nil, err
	}
	if err != nil {
		logger.Fatalf("Failed to generate private key: %v", err)
		return nil, err
	}

	var notBefore time.Time
	if len(gc.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", gc.ValidFrom)
		if err != nil {
			logger.Fatalf("Failed to parse creation date: %v", err)
			return nil, err
		}
	}

	template, err := gc.createTemplate(notBefore)
	if err != nil {
		logger.Fatalf("Failed to create certificate template: %v", err)
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, gc.publicKey(priv), priv)
	if err != nil {
		logger.Fatalf("Failed to create certificate: %v", err)
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		logger.Fatalf("Unable to marshal private key: %v", err)
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		logger.Fatalf("Failed to combine tls key pair: %v", err)
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}

func (gc *generateCertificate) createTemplate(notBefore time.Time) (*x509.Certificate, error) {
	notAfter := notBefore.Add(gc.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Fatalf("Failed to generate serial number: %v", err)
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

// GetTLSConfig generates the tls config, getting certificate from ca service
func (s *SHttp) GetTLSConfig() (*tls.Config, error) {
	var err error
	subj := pkix.Name{
		CommonName: s.cfa.Servicename,
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

	for _, sip := range s.cfn.IPAddresses {
		if ip := net.ParseIP(sip); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		}
	}

	for _, sdn := range s.cfn.DNSNames {
		template.DNSNames = append(template.DNSNames, sdn)
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

	cli, err := mv.LoginClient(s.cfa.AccessKey, s.cfa.Secret, s.cfa.URL)
	if err != nil {
		return nil, err
	}

	crt, err := cli.CreateCertificate(template)
	if err != nil {
		return nil, err
	}

	prv, err := cli.PrivateKey()
	if err != nil {
		return nil, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(prv)
	if err != nil {
		logger.Fatalf("Unable to marshal private key: %v", err)
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

func (s *SHttp) TLSFromFiles() (*tls.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair(s.cfn.Certificate, s.cfn.Key)
	if err != nil {
		return nil, errors.Wrap(err, "error creating X509 Pair")
	}

	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}, nil
}
