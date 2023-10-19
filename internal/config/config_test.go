package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

const (
	testDataPath     = "./../../testdata/"
	serviceLocalFile = "service_local_file.yaml"
	logFileName      = "file.log"
)

func TestLoadFromYaml(t *testing.T) {
	ast := assert.New(t)
	File = testDataPath + serviceLocalFile

	err := Load()
	ast.Nil(err)
	c := Get()

	ast.Equal(8000, Get().Services.HTTP.Port)
	ast.Equal(8443, Get().Services.HTTP.Sslport)

	ast.Equal(60, Get().Services.HealthSystem.Period)
	ast.Equal(3, Get().Services.HealthSystem.StartDelay)
	ast.Equal("", Get().SecretFile)
	ast.Equal("https://127.0.0.1:8443", Get().Services.HTTP.ServiceURL)
	c.Provide()

	cfg := do.MustInvoke[Config](nil)
	ast.Nil(err)
	ast.NotNil(cfg)
	do.MustShutdown[Config](nil)
}

func TestDefaultConfig(t *testing.T) {
	ast := assert.New(t)
	config = DefaultConfig
	cnf := Get()

	ast.Equal(8000, cnf.Services.HTTP.Port)
	ast.Equal(8443, cnf.Services.HTTP.Sslport)

	ast.Equal(30, cnf.Services.HealthSystem.Period)
	ast.Equal(3, cnf.Services.HealthSystem.StartDelay)
	ast.Equal("", cnf.SecretFile)
	ast.Equal("https://127.0.0.1:8443", cnf.Services.HTTP.ServiceURL)

	ast.Equal("INFO", cnf.Logging.Level)
}

func TestCfgSubst(t *testing.T) {
	ast := assert.New(t)

	File = filepath.Join("${configdir}", serviceLocalFile)

	err := Load()
	ast.NotNil(err)
	home, err := os.UserConfigDir()
	ast.Nil(err)
	file := filepath.Join(home, Servicename, serviceLocalFile)
	ast.Equal(file, File)
}

func TestEnvSubstRightCase(t *testing.T) {
	ast := assert.New(t)

	err := os.Setenv("logfile", logFileName)
	ast.Nil(err)

	File = testDataPath + serviceLocalFile

	err = Load()
	ast.Nil(err)

	ast.Equal(logFileName, Get().Logging.Filename)
}

func TestEnvSubstWrongCase(t *testing.T) {
	ast := assert.New(t)

	err := os.Setenv("LogFile", logFileName)
	ast.Nil(err)

	File = testDataPath + serviceLocalFile

	err = Load()
	ast.Nil(err)

	ast.Equal(logFileName, Get().Logging.Filename)
}

func TestSecretMapping(t *testing.T) {
	ast := assert.New(t)

	File = "./../../testdata/service_local_file_w_secret.yaml"

	err := Load()
	ast.Nil(err)

	ast.Equal(120, Get().Services.HealthSystem.Period)
}
