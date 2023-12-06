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
	cl, err = NewClient("https://127.0.0.1:9443", Tenant)
	if err != nil {
		panic(err)
	}
}

func TestClientCRUD(t *testing.T) {
	initCl()
	ast := assert.New(t)

	cd := pmodel.Address{
		Name:      "Smith",
		Firstname: "John",
		Street:    "123 Main St",
		City:      "Anytown",
		State:     "CA",
		ZipCode:   "12345",
	}

	id, err := cl.CreateAddress(cd)
	ast.Nil(err)
	ast.NotEmpty(id)

	cd1, err := cl.GetAddress(id)
	ast.Nil(err)
	ast.NotNil(cd1)

	ast.Equal(id, cd1.ID)
	ast.Equal(cd.Name, cd1.Name)
	ast.Equal(cd.Firstname, cd1.Firstname)
	ast.Equal(cd.City, cd1.City)
	ast.Equal(cd.State, cd1.State)
	ast.Equal(cd.Street, cd1.Street)
	ast.Equal(cd.ZipCode, cd1.ZipCode)

	l, err := cl.GetAddresses()
	ast.Nil(err)
	ast.Equal(1, len(*l))

	ok, err := cl.DeleteAddress(id)
	ast.Nil(err)
	ast.True(ok)

	cd1, err = cl.GetAddress(id)
	ast.NotNil(err)
	ast.Nil(cd1)

	ok, err = cl.DeleteAddress(id)
	ast.Nil(err)
	ast.False(ok)
}
