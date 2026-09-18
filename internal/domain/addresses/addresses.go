package addresses

import (
	"context"

	"github.com/willie68/go-micro/pkg/pmodel"
)

// This is the domain model for CRUD Operations to address storage
// AddressStorage is the outbound port for address persistence.
type AddressStorage interface {
	Addresses() ([]pmodel.Address, error)
	Has(id string) bool
	Read(id string) (*pmodel.Address, error)
	Create(adr pmodel.Address) (string, error)
	Update(adr pmodel.Address) error
	Delete(id string) error
}

type Addresses struct {
	storage AddressStorage
}

// NewAddresses creates a new Addresses instance
func NewAddresses(storage AddressStorage) *Addresses {
	return &Addresses{
		storage: storage,
	}
}

// Addresses returns all addresses
func (a *Addresses) Addresses(ctx context.Context) ([]pmodel.Address, error) {
	return a.storage.Addresses()
}

// Has checks if an address exists
func (a *Addresses) Has(ctx context.Context, id string) bool {
	return a.storage.Has(id)
}

// Read reads an address
func (a *Addresses) Read(ctx context.Context, id string) (*pmodel.Address, error) {
	return a.storage.Read(id)
}

// Create creates a new address
func (a *Addresses) Create(ctx context.Context, adr pmodel.Address) (string, error) {
	return a.storage.Create(adr)
}

// Update updates an address
func (a *Addresses) Update(ctx context.Context, adr pmodel.Address) error {
	return a.storage.Update(adr)
}

// Delete deletes an address
func (a *Addresses) Delete(ctx context.Context, id string) error {
	return a.storage.Delete(id)
}
