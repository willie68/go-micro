package adrint

import (

	// needed declaration
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/pkg/pmodel"
)

// AdrInt the internal address storage type
type AdrInt struct {
	adrs map[string]pmodel.Address
}

var _ adrsvc.AddressStorage = &AdrInt{}

// NewAdrInt create a new instance of the internal address storage
func NewAdrInt() (*AdrInt, error) {
	am := AdrInt{
		adrs: make(map[string]pmodel.Address),
	}
	return &am, nil
}

// Addresses list all addresses
func (a *AdrInt) Addresses() ([]pmodel.Address, error) {
	addresses := make([]pmodel.Address, 0)
	for _, v := range a.adrs {
		addresses = append(addresses, v)
	}
	return addresses, nil
}

// Has checking if an adress is present
func (a *AdrInt) Has(id string) bool {
	_, ok := a.adrs[id]
	return ok
}

// Read getting the address with id
func (a *AdrInt) Read(id string) (*pmodel.Address, error) {
	adr, ok := a.adrs[id]
	if !ok {
		return nil, adrsvc.ErrNotFound
	}
	return &adr, nil
}

// Create creates a new Address
func (a *AdrInt) Create(adr pmodel.Address) (string, error) {
	id := xid.New().String()
	adr.ID = id
	a.adrs[id] = adr
	return id, nil
}

// Update updates the address
func (a *AdrInt) Update(adr pmodel.Address) error {
	_, ok := a.adrs[adr.ID]
	if !ok {
		return adrsvc.ErrNotFound
	}
	a.adrs[adr.ID] = adr
	return nil
}

// Delete deletes the address with id
func (a *AdrInt) Delete(id string) error {
	_, ok := a.adrs[id]
	if !ok {
		return adrsvc.ErrNotFound
	}
	delete(a.adrs, id)
	return nil
}

// CheckName should return the name of this healthcheck. The name should be unique.
func (a *AdrInt) CheckName() string {
	return "internal"
}

// Check proceed a check and return state, true for healthy or false and an optional error, if the healthcheck fails
func (a *AdrInt) Check() (bool, error) {
	return true, nil
}

// Init this service
func (a *AdrInt) Init() error {
	// Nothing to do here
	return nil
}

// Shutdown this service
func (a *AdrInt) Shutdown() error {
	// Nothing to do here
	return nil
}
