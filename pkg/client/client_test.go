package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/pkg/pmodel"
)

const (
	Tenant = "tester1"
)

var cl *Client

func initCl() {
	var err error
	StartServer()
	cl, err = NewClient("https://127.0.0.1:8443", Tenant)
	if err != nil {
		panic(err)
	}
}

func TestClientCRUD(t *testing.T) {
	initCl()
	ast := assert.New(t)

	cd := pmodel.ConfigDescription{
		StoreID:  "myNewStore",
		TenantID: Tenant,
		Size:     1234567,
	}

	id, err := cl.PutConfig(cd)
	ast.Nil(err)
	ast.Equal(cd.TenantID, id)

	cd1, err := cl.GetConfig(cd.TenantID)
	ast.Nil(err)
	ast.NotNil(cd1)

	ast.Equal("myNewStore", cd1.StoreID)
	ast.Equal(1234567, cd1.Size)
	ast.Equal(Tenant, cd1.TenantID)

	cd1, err = cl.GetMyConfig()
	ast.Nil(err)
	ast.NotNil(cd1)

	ast.Equal("myNewStore", cd1.StoreID)
	ast.Equal(1234567, cd1.Size)
	ast.Equal(Tenant, cd1.TenantID)

	l, err := cl.GetConfigs()
	ast.Nil(err)
	ast.Equal(1, len(*l))

	ok, err := cl.DeleteConfig(cd.TenantID)
	ast.Nil(err)
	ast.True(ok)

	cd1, err = cl.GetConfig(cd.TenantID)
	ast.NotNil(err)
	ast.Nil(cd1)

	ok, err = cl.DeleteConfig(cd.TenantID)
	ast.Nil(err)
	ast.False(ok)
}
