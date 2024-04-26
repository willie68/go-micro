package adrmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	// needed declaration
	_ "github.com/go-sql-driver/mysql"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/pkg/pmodel"
)

// Config configuration for mysql
type Config struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// AdrMdb this is the address mysql type
type AdrMdb struct {
	db   *sql.DB
	mcfg Config
}

var _ adrsvc.AddressStorage = &AdrMdb{}

// NewAdrMdb creates a new AdrMdb isntance
func NewAdrMdb(cfg Config) (*AdrMdb, error) {
	d, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Database))
	if err != nil {
		return nil, err
	}
	am := AdrMdb{
		db:   d,
		mcfg: cfg,
	}
	return &am, nil
}

// Addresses list all addresses
func (a *AdrMdb) Addresses() ([]pmodel.Address, error) {
	addresses := []pmodel.Address{}

	rows, err := a.db.Query(fmt.Sprintf("SELECT id, name, firstname, street, city, state, zip_code FROM %s", a.mcfg.Table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var address pmodel.Address
		err := rows.Scan(&address.ID, &address.Name, &address.Firstname, &address.Street, &address.City, &address.State, &address.ZipCode)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

// Has checking if an adress is present
func (a *AdrMdb) Has(id string) bool {
	var mid string
	err := a.db.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE id=?", a.mcfg.Table), id).Scan(&mid)
	return err == nil
}

// Read getting the address with id
func (a *AdrMdb) Read(id string) (*pmodel.Address, error) {
	var address pmodel.Address

	err := a.db.QueryRow(fmt.Sprintf("SELECT id, name, firstname, street, city, state, zip_code FROM %s WHERE id=?", a.mcfg.Table), id).Scan(
		&address.ID, &address.Name, &address.Firstname, &address.Street, &address.City, &address.State, &address.ZipCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, adrsvc.ErrNotFound
		}
		return nil, err
	}
	return &address, nil
}

// Create creates a new Address
func (a *AdrMdb) Create(adr pmodel.Address) (string, error) {
	result, err := a.db.Exec(fmt.Sprintf("INSERT INTO %s (name, firstname, street, city, state, zip_code) VALUES (?, ?, ?, ?, ?, ?)",
		a.mcfg.Table), adr.Name, adr.Firstname, adr.Street, adr.City, adr.State, adr.ZipCode)

	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", nil
	}

	return strconv.FormatInt(id, 10), nil
}

// Update updates the address
func (a *AdrMdb) Update(adr pmodel.Address) error {
	_, err := a.db.Exec(fmt.Sprintf("UPDATE %s SET name=?, firstname=?, street=?, city=?, state=?, zip_code=? WHERE id=?", a.mcfg.Table), adr.Name, adr.Firstname, adr.Street, adr.City, adr.State, adr.ZipCode, adr.ID)

	return err
}

// Delete deletes the address with id
func (a *AdrMdb) Delete(id string) error {
	_, err := a.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=?", a.mcfg.Table), id)
	return err
}

// CheckName should return the name of this healthcheck. The name should be unique.
func (a *AdrMdb) CheckName() string {
	return "mysql"
}

// Check proceed a check and return state, true for healthy or false and an optional error, if the healthcheck fails
func (a *AdrMdb) Check() (bool, error) {
	err := a.db.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Init this service
func (a *AdrMdb) Init() error {
	// Nothing to do here
	return nil
}

// Shutdown this service
func (a *AdrMdb) Shutdown() error {
	// Nothing to do here
	return nil
}
