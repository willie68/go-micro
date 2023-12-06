package adrmysql

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/pkg/pmodel"
)

var adrs []pmodel.Address

func init() {
	data, err := os.ReadFile("../../../../testdata/addresses.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(
		data, &adrs)
	if err != nil {
		panic(err)
	}
}

func TestAdrMdbList(t *testing.T) {
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sdb.Close()

	cfg := Config{
		Host:     "127.0.0.1",
		Database: "golang",
		Table:    "address",
		Username: "address",
		Password: "address",
	}
	stg := AdrMdb{
		db:   sdb,
		mcfg: cfg,
	}
	ast.NotNil(stg)

	// List
	rows := sqlmock.NewRows([]string{"id", "lastname", "firstname", "street", "city", "state", "zip_code"})
	for _, adr := range adrs {
		str := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s ", adr.ID, adr.Name, adr.Firstname, adr.Street, adr.City, adr.State, adr.ZipCode)
		rows = rows.FromCSVString(str)
	}
	mock.ExpectQuery("SELECT id, name, firstname, street, city, state, zip_code FROM address").WillReturnRows(rows)

	as, err := stg.Addresses()
	ast.Nil(err)
	ast.NotNil(as)
	ast.Len(as, 10)
}

func TestAdrMdbRead(t *testing.T) {
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sdb.Close()

	cfg := Config{
		Host:     "127.0.0.1",
		Database: "golang",
		Table:    "address",
		Username: "address",
		Password: "address",
	}
	stg := AdrMdb{
		db:   sdb,
		mcfg: cfg,
	}
	ast.NotNil(stg)

	// List
	rows := sqlmock.NewRows([]string{"id", "lastname", "firstname", "street", "city", "state", "zip_code"})
	adr := adrs[3]
	str := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s ", adr.ID, adr.Name, adr.Firstname, adr.Street, adr.City, adr.State, adr.ZipCode)
	rows = rows.FromCSVString(str)
	mock.ExpectQuery("SELECT id, name, firstname, street, city, state, zip_code FROM address WHERE id=?").WillReturnRows(rows)

	as, err := stg.Read("4")
	ast.Nil(err)
	ast.NotNil(as)
	ast.Equal("4", as.ID)
}

func TestAdrMdbCreate(t *testing.T) {
	t.SkipNow()
}

func TestAdrMdbUpdate(t *testing.T) {
	t.SkipNow()
}
func TestAdrMdbDelete(t *testing.T) {
	t.SkipNow()
}
