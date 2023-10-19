package sconfig

import (
	"errors"
	"fmt"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/model"
	"github.com/willie68/go-micro/internal/serror"
)

// TestDI checking if SConfig is registered in the DI Framework
func TestDI(t *testing.T) {
	ast := assert.New(t)
	_, err := NewSConfig()
	ast.Nil(err)

	sc2 := do.MustInvoke[*SConfig](nil)
	ast.NotNil(sc2)

	do.MustShutdown[SConfig](nil)
}

func TestCRUD(t *testing.T) {
	ast := assert.New(t)
	sc, err := NewSConfig()
	ast.Nil(err)

	cd := model.ConfigDescription{
		StoreID:  "12345678",
		TenantID: "mytenant",
		Size:     12345,
	}

	ast.False(sc.HasConfig(cd.TenantID))

	n, err := sc.PutConfig(cd.TenantID, cd)
	ast.Nil(err)
	ast.Equal(cd.TenantID, n)

	ast.True(sc.HasConfig(cd.TenantID))

	cd2, err := sc.GetConfig(cd.TenantID)
	ast.Nil(err)
	ast.NotNil(cd2)

	ast.Equal(cd.StoreID, cd2.StoreID)
	ast.Equal(cd.TenantID, cd2.TenantID)
	ast.Equal(cd.Size, cd2.Size)

	ast.True(sc.DeleteConfig(cd.TenantID))

	ast.False(sc.HasConfig(cd.TenantID))

	cd2, err = sc.GetConfig(cd.TenantID)
	ast.NotNil(err)
	ast.Nil(cd2)

	do.MustShutdown[SConfig](nil)
}

func TestNIY(t *testing.T) {
	ast := assert.New(t)
	sc, err := NewSConfig()
	ast.Nil(err)

	err = sc.NotImplemented()
	ast.NotNil(err)
	ast.True(errors.Is(serror.ErrNotImplementedYet, err))

	do.MustShutdown[SConfig](nil)
}

func TestList(t *testing.T) {
	ast := assert.New(t)
	sc, err := NewSConfig()
	ast.Nil(err)

	for i := 1; i < 100; i++ {
		cd := model.ConfigDescription{
			StoreID:  "12345678",
			TenantID: fmt.Sprintf("mytenant%d", i),
			Size:     i,
		}
		_, err := sc.PutConfig(cd.TenantID, cd)
		ast.Nil(err)
	}

	l, err := sc.ListConfigs()
	ast.Nil(err)
	ast.Equal(99, len(l))

	do.MustShutdown[SConfig](nil)
}
