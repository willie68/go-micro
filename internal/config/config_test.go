package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

const (
	testDataPath     = "./../../testdata/"
	serviceLocalFile = "service_local_file.yaml"
	logFileName      = "file.log"
)

func TestLoadFromYaml(t *testing.T) {
	inj := do.New()
	ast := assert.New(t)
	File = testDataPath + serviceLocalFile

	err := Load()
	ast.Nil(err)
	c := Get()

	ast.Equal(8000, Get().HTTP.Port)
	ast.Equal(8443, Get().HTTP.Sslport)

	ast.Equal(60, Get().HealthSystem.Period)
	ast.Equal(3, Get().HealthSystem.StartDelay)
	ast.Equal("", Get().SecretFile)
	ast.Equal("https://127.0.0.1:8443", Get().HTTP.ServiceURL)
	c.Provide(inj)

	cfg := do.MustInvoke[Config](inj)
	ast.Nil(err)
	ast.NotNil(cfg)
	do.MustShutdown[Config](inj)
}

func TestDefaultConfig(t *testing.T) {
	ast := assert.New(t)
	config = DefaultConfig
	cnf := Get()

	ast.Equal(8000, cnf.HTTP.Port)
	ast.Equal(8443, cnf.HTTP.Sslport)

	ast.Equal(30, cnf.HealthSystem.Period)
	ast.Equal(3, cnf.HealthSystem.StartDelay)
	ast.Equal("", cnf.SecretFile)
	ast.Equal("https://127.0.0.1:8443", cnf.HTTP.ServiceURL)

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

	ast.Equal(120, Get().HealthSystem.Period)
}
