package config

// Config our service configuration
type Config struct {
	//port of the http server
	Port int `yaml:"port"`
	//port of the https server
	Sslport int `yaml:"sslport"`
	//this is the url how to connect to this service from outside
	ServiceURL string `yaml:"serviceURL"`
	//this is the url where to register this service
	RegistryURL string `yaml:"registryURL"`
	//this is the url where to register this service
	SystemID string `yaml:"systemID"`

	SecretFile string  `yaml:"secretfile"`
	Logging    Logging `yaml:"logging"`

	HealthCheck HealthCheck `yaml:"healthcheck"`

	OpenTracing OpenTracing `yaml:"opentracing"`
}

// Logging configuration for the gelf logging
type Logging struct {
	Gelfurl  string `yaml:"gelf-url"`
	Gelfport int    `yaml:"gelf-port"`
}

// HealthCheck configuration for the health check system
type HealthCheck struct {
	Period int `yaml:"period"`
}

type OpenTracing struct {
	Host     string `yaml:"host"`
	Endpoint string `yaml:"endpoint"`
}
