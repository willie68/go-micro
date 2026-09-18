package adrint

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/pkg/pmodel"
)

var (
	madrs map[string]pmodel.Address
	adrs  []pmodel.Address
)

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
	madrs = make(map[string]pmodel.Address)
	for _, v := range adrs {
		madrs[v.ID] = v
	}
}

func TestAdrMdbList(t *testing.T) {
	ast := assert.New(t)
	stg := AdrInt{
		adrs: madrs,
	}
	ast.NotNil(stg)

	as, err := stg.Addresses()
	ast.Nil(err)
	ast.NotNil(as)
	ast.Len(as, 10)
}

func TestAdrMdbRead(t *testing.T) {
	ast := assert.New(t)
	stg := AdrInt{
		adrs: madrs,
	}
	ast.NotNil(stg)

	as, err := stg.Read("4")
	ast.Nil(err)
	ast.NotNil(as)
	ast.Equal("4", as.ID)
}

func TestAdrMdbCreate(t *testing.T) {
	ast := assert.New(t)
	stg, err := NewAdrInt()
	ast.NoError(err)

	id, err := stg.Create(pmodel.Address{Name: "Doe", Firstname: "John"})
	ast.NoError(err)
	ast.NotEmpty(id)
	ast.True(stg.Has(id))

	got, err := stg.Read(id)
	ast.NoError(err)
	ast.Equal("Doe", got.Name)
}

func TestAdrMdbUpdate(t *testing.T) {
	ast := assert.New(t)
	stg := AdrInt{adrs: map[string]pmodel.Address{"4": madrs["4"]}}

	err := stg.Update(pmodel.Address{ID: "missing"})
	ast.Error(err)

	updated := madrs["4"]
	updated.Name = "Changed"
	ast.NoError(stg.Update(updated))
	got, err := stg.Read("4")
	ast.NoError(err)
	ast.Equal("Changed", got.Name)
}

func TestAdrMdbDelete(t *testing.T) {
	ast := assert.New(t)
	stg := AdrInt{adrs: map[string]pmodel.Address{"4": madrs["4"]}}

	ast.Error(stg.Delete("missing"))
	ast.NoError(stg.Delete("4"))
	ast.False(stg.Has("4"))
	_, err := stg.Read("4")
	ast.Error(err)
}

func TestAdrIntHealthAndLifecycle(t *testing.T) {
	stg, err := NewAdrInt()
	assert.NoError(t, err)
	assert.Equal(t, "internal", stg.CheckName())
	ok, err := stg.Check()
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.NoError(t, stg.Init())
	assert.NoError(t, stg.Shutdown())
	assert.NoError(t, stg.HealthCheck())
}
