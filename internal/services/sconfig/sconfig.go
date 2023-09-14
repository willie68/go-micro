package sconfig

import (
	"errors"
	"sync"

	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/model"
	"github.com/willie68/go-micro/internal/serror"
)

// DoConfig dependency injection key name for sconfig
const DoConfig = "configs"

// SConfig config management
type SConfig struct {
	cfgs sync.Map

	healthy bool
}

// NewSConfig creating a new config business object
func NewSConfig() (*SConfig, error) {
	gs := SConfig{
		cfgs: sync.Map{},
	}
	do.ProvideNamedValue[*SConfig](nil, DoConfig, &gs)
	return &gs, nil
}

// Init initialize the config management
func (s *SConfig) Init() error {
	return nil
}

// GetConfig adding a new config to the service
func (s *SConfig) GetConfig(n string) (*model.ConfigDescription, error) {
	a, ok := s.cfgs.Load(n)
	if !ok {
		return nil, serror.ErrNotExists
	}
	cd, ok := a.(model.ConfigDescription)
	if !ok {
		return nil, serror.ErrUnknowError
	}
	return &cd, nil
}

// PutConfig adding/Updating a config
func (s *SConfig) PutConfig(n string, c model.ConfigDescription) (id string, err error) {
	s.cfgs.Store(n, c)
	return n, nil
}

// DeleteConfig deleting a config
func (s *SConfig) DeleteConfig(n string) bool {
	_, ok := s.cfgs.LoadAndDelete(n)
	return ok
}

// HasConfig checking if a config is already added
func (s *SConfig) HasConfig(n string) bool {
	_, ok := s.cfgs.Load(n)
	return ok
}

// List list all names of configs
func (s *SConfig) List() ([]string, error) {
	l := make([]string, 0)
	s.cfgs.Range(func(key, value any) bool {
		n, ok := key.(string)
		if ok {
			l = append(l, n)
		}
		return true
	})
	return l, nil
}

// Name needed for the health check system
func (s *SConfig) CheckName() string {
	return "sconfig-service"
}

// Check proceed a check and return state, true for healthy or false and an optional error, if the healthcheck fails
func (s *SConfig) Check() (bool, error) {
	s.healthy = !s.healthy
	if s.healthy {
		return true, nil
	} else {
		return false, errors.New("sconfig not healthy")
	}
}

// NotImplemented throwing a not implemented error
func (s *SConfig) NotImplemented() error {
	return serror.ErrNotImplementedYet
}
