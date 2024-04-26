package shttp

// Config configuration of the http service
type Config struct {
	// port of the http server
	Port int `yaml:"port"`
	// port of the https server
	Sslport int `yaml:"sslport"`
	// this is the url how to connect to this service from outside
	ServiceURL string `yaml:"serviceURL"`
	// other dns names (used for certificate)
	DNSNames []string `yaml:"dnss"`
	// other ips (used for certificate)
	IPAddresses []string `yaml:"ips"`
	// Servicename for the certificate
	Servicename string `yaml:"servicename"`
	// path and name to the certificate
	Certificate string `yaml:"certificate"`
	// path and name to the private key
	Key string `yaml:"key"`
}
