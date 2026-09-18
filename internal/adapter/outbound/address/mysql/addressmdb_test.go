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
	data, err := os.ReadFile("../../../../../testdata/addresses.json")
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
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New()
	ast.NoError(err)
	defer sdb.Close()

	stg := AdrMdb{db: sdb, mcfg: Config{Table: "address"}}
	mock.ExpectExec("INSERT INTO address").WillReturnResult(sqlmock.NewResult(42, 1))

	id, err := stg.Create(pmodel.Address{Name: "Doe"})
	ast.NoError(err)
	ast.Equal("42", id)
}

func TestAdrMdbUpdate(t *testing.T) {
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New()
	ast.NoError(err)
	defer sdb.Close()

	stg := AdrMdb{db: sdb, mcfg: Config{Table: "address"}}
	mock.ExpectExec("UPDATE address SET").WillReturnResult(sqlmock.NewResult(0, 1))
	ast.NoError(stg.Update(pmodel.Address{ID: "4", Name: "Doe"}))
}

func TestAdrMdbDelete(t *testing.T) {
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New()
	ast.NoError(err)
	defer sdb.Close()

	stg := AdrMdb{db: sdb, mcfg: Config{Table: "address"}}
	mock.ExpectExec("DELETE FROM address").WillReturnResult(sqlmock.NewResult(0, 1))
	ast.NoError(stg.Delete("4"))
}

func TestAdrMdbHasAndHealth(t *testing.T) {
	ast := assert.New(t)
	sdb, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	ast.NoError(err)
	defer sdb.Close()

	stg := AdrMdb{db: sdb, mcfg: Config{Table: "address"}}
	mock.ExpectQuery("SELECT id FROM address WHERE id=?").WithArgs("4").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("4"))
	ast.True(stg.Has("4"))

	mock.ExpectQuery("SELECT id FROM address WHERE id=?").WithArgs("x").
		WillReturnError(fmt.Errorf("sql: no rows"))
	ast.False(stg.Has("x"))

	ast.Equal("mysql", stg.CheckName())
	mock.ExpectPing()
	ok, err := stg.Check()
	ast.NoError(err)
	ast.True(ok)
	ast.NoError(stg.Init())
	ast.NoError(stg.Shutdown())
}
