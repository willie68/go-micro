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
	data, err := os.ReadFile("../../../../testdata/addresses.json")
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
	t.SkipNow()
}

func TestAdrMdbUpdate(t *testing.T) {
	t.SkipNow()
}
func TestAdrMdbDelete(t *testing.T) {
	t.SkipNow()
}
