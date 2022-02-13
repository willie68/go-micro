package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromYaml(t *testing.T) {
	ast := assert.New(t)
	File = "./../../testdata/service_local_file.yaml"

	Load()

	ast.Equal(8000, Get().Port)
	ast.Equal(8443, Get().Sslport)

	ast.Equal(60, Get().HealthCheck.Period)
	ast.Equal("", Get().SecretFile)
	ast.Equal("https://127.0.0.1:8443", Get().ServiceURL)
}

func TestDefaultConfig(t *testing.T) {
	ast := assert.New(t)
	config = DefaultConfig
	cnf := Get()

	ast.Equal(8000, cnf.Port)
	ast.Equal(8443, cnf.Sslport)

	ast.Equal(30, cnf.HealthCheck.Period)
	ast.Equal("", cnf.SecretFile)
	ast.Equal("https://127.0.0.1:8443", cnf.ServiceURL)

	ast.Equal("INFO", cnf.Logging.Level)
}

func TestCfgSubst(t *testing.T) {
	ast := assert.New(t)

	File = filepath.Join("${configdir}", "service_local_file.yaml")

	err := Load()
	ast.NotNil(err)
	home, err := os.UserConfigDir()
	ast.Nil(err)
	file := filepath.Join(home, Servicename, "service_local_file.yaml")
	ast.Equal(file, File)
}

func TestEnvSubstRightCase(t *testing.T) {
	ast := assert.New(t)

	os.Setenv("logfile", "file.log")

	File = "./../../testdata/service_local_file.yaml"

	Load()

	ast.Equal("file.log", Get().Logging.Filename)
}

func TestEnvSubstWrongCase(t *testing.T) {
	ast := assert.New(t)

	os.Setenv("LogFile", "file.log")

	File = "./../../testdata/service_local_file.yaml"

	Load()

	ast.Equal("file.log", Get().Logging.Filename)
}
