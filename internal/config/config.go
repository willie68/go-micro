package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const Servicename = "go-micro"

// Config our service configuration
type Config struct {
	//port of the http server
	Port int `yaml:"port"`
	//port of the https server
	Sslport int `yaml:"sslport"`
	//this is the url how to connect to this service from outside
	ServiceURL string `yaml:"serviceURL"`

	SecretFile string `yaml:"secretfile"`
	Apikey     bool   `yaml:"apikey"`

	Logging LoggingConfig `yaml:"logging"`

	HealthCheck HealthCheck `yaml:"healthcheck"`

	Auth Authentcation `yaml:"auth"`

	OpenTracing OpenTracing `yaml:"opentracing"`

	Metrics Metrics `yaml:"metrics"`
}

type Authentcation struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
}

// HealthCheck configuration for the health check system
type HealthCheck struct {
	Period int `yaml:"period"`
}

// Logging configuration for the gelf logging
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`

	Gelfurl  string `yaml:"gelf-url"`
	Gelfport int    `yaml:"gelf-port"`
}

type OpenTracing struct {
	Host     string `yaml:"host"`
	Endpoint string `yaml:"endpoint"`
}

type Metrics struct {
	Enable bool `yaml:"enable"`
}

var DefaultConfig = Config{
	Port:       8000,
	Sslport:    8443,
	ServiceURL: "https://127.0.0.1:8443",
	SecretFile: "",
	Apikey:     true,
	HealthCheck: HealthCheck{
		Period: 30,
	},
	Logging: LoggingConfig{
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
	data, err := ioutil.ReadFile(File)
	if err != nil {
		return fmt.Errorf("can't load config file: %s", err.Error())
	}
	dataStr := os.ExpandEnv(string(data))
	err = yaml.Unmarshal([]byte(dataStr), &config)
	if err != nil {
		return fmt.Errorf("can't unmarshal config file: %s", err.Error())
	}
	return readSecret()
}

func readSecret() error {
	secretFile := config.SecretFile
	if secretFile != "" {
		data, err := ioutil.ReadFile(secretFile)
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
