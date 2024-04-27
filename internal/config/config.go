package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/drone/envsubst"
	"github.com/pkg/errors"
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/internal/services/caservice"
	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/shttp"
	"gopkg.in/yaml.v3"
)

// Servicename the name of this service
const Servicename = "go-micro"

// Config our service configuration
type Config struct {
	// all secrets will be stored in this file, same structure as the main config file
	SecretFile string `yaml:"secretfile"`
	// configure logging to gelf logging system
	Logging logging.Config `yaml:"logging"`
	// use authentication via jwt
	Auth Authentication `yaml:"auth"`
	// opentelemtrie tracer can be configured here
	OpenTracing OpenTracing `yaml:"opentracing"`
	// and some metrics
	Metrics Metrics `yaml:"metrics"`
	// HTTP REST Service
	HTTP shttp.Config `yaml:"http"`
	// special config for health checks
	HealthSystem health.Config `yaml:"healthcheck"`
	// CA service will be used, microvault
	CA caservice.Config `yaml:"ca"`
	// Enable Profiling option
	Profiling Profiling `yaml:"profiling"`
	// This is the demo address storage config
	AddressStorage adrsvc.Config `yaml:"addressstorage"`
}

// Authentication configuration
type Authentication struct {
	Type       string         `yaml:"type"`
	Properties map[string]any `yaml:"properties"`
}

// OpenTracing configuration
type OpenTracing struct {
	Host     string `yaml:"host"`
	Endpoint string `yaml:"endpoint"`
}

// Metrics configuration
type Metrics struct {
	Enable bool `yaml:"enable"`
}

// Profiling configuration
type Profiling struct {
	Enable bool `yaml:"enable"`
}

// DefaultConfig default configuration
var DefaultConfig = Config{
	HTTP: shttp.Config{
		Servicename: "go-micro",
		Port:        8000,
		Sslport:     8443,
		ServiceURL:  "https://127.0.0.1:8443",
	},
	HealthSystem: health.Config{
		Period:     30,
		StartDelay: 3,
	},
	SecretFile: "",
	Logging: logging.Config{
		Level:    "INFO",
		Filename: "${configdir}/logging.log",
	},
}

// GetDefaultConfigFolder returning the default configuration folder of the system
func GetDefaultConfigFolder() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFolder := filepath.Join(home, Servicename)
	err = os.MkdirAll(configFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	return configFolder, nil
}

// GetDefaultConfigfile getting the default config file
func GetDefaultConfigfile() (string, error) {
	configFolder, err := GetDefaultConfigFolder()
	if err != nil {
		return "", errors.Wrap(err, "can't load config file")
	}
	configFolder = filepath.Join(configFolder, "service")
	err = os.MkdirAll(configFolder, os.ModePerm)
	if err != nil {
		return "", errors.Wrap(err, "can't load config file")
	}
	return filepath.Join(configFolder, "service.yaml"), nil
}

// ReplaceConfigdir replace the configdir macro
func ReplaceConfigdir(s string) (string, error) {
	if strings.Contains(s, "${configdir}") {
		configFolder, err := GetDefaultConfigFolder()
		if err != nil {
			return "", err
		}
		return strings.Replace(s, "${configdir}", configFolder, -1), nil
	}
	return s, nil
}

var config = Config{}

// File the config file
var File = "${configdir}/service.yaml"

func init() {
	config = DefaultConfig
}

// Provide provide the config to the dependency injection
func (c *Config) Provide() {
	do.ProvideValue[Config](nil, *c)
}

// Get returns loaded config
func Get() Config {
	return config
}

// Load loads the config
func Load() error {
	myFile, err := ReplaceConfigdir(File)
	if err != nil {
		return fmt.Errorf("can't get default config folder: %s", err.Error())
	}
	File = myFile
	_, err = os.Stat(myFile)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(File)
	if err != nil {
		return fmt.Errorf("can't load config file: %s", err.Error())
	}
	dataStr, err := envsubst.EvalEnv(string(data))
	if err != nil {
		return fmt.Errorf("can't substitute config file: %s", err.Error())
	}

	err = yaml.Unmarshal([]byte(dataStr), &config)
	if err != nil {
		return fmt.Errorf("can't unmarshal config file: %s", err.Error())
	}
	return readSecret()
}

func readSecret() error {
	secretFile := config.SecretFile
	if secretFile != "" {
		data, err := os.ReadFile(secretFile)
		if err != nil {
			return fmt.Errorf("can't load secret file: %s", err.Error())
		}
		var secretConfig Config
		err = yaml.Unmarshal(data, &secretConfig)
		if err != nil {
			return fmt.Errorf("can't unmarshal secret file: %s", err.Error())
		}
		// merge secret
		if err := mergo.Map(&config, secretConfig, mergo.WithOverride); err != nil {
			return fmt.Errorf("can't merge secret file: %s", err.Error())
		}
	}
	return nil
}
