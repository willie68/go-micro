package adrsvc

import (
	"errors"

	"github.com/willie68/go-micro/internal/services/health"
	"github.com/willie68/go-micro/internal/services/interfaces"
	"github.com/willie68/go-micro/pkg/pmodel"
)

// Error definitions
var (
	ErrNotFound = errors.New("not found")
)

// Config general config for the storage
type Config struct {
	Type       string         `yaml:"type"`
	Connection map[string]any `yaml:"connection"`
}

// AddressStorage this is the interface every storage provider should implement
//
//go:generate mockery --name=AddressStorage --outpkg=mocks --filename=addressstorage.go --with-expecter
type AddressStorage interface {
	interfaces.Service
	health.Check
	Addresses() ([]pmodel.Address, error)
	Has(id string) bool
	Read(id string) (*pmodel.Address, error)
	Create(adr pmodel.Address) (string, error)
	Update(adr pmodel.Address) error
	Delete(id string) error
}
