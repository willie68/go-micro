package adrfact

import (
	"github.com/samber/do"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/internal/services/adrsvc/adrint"
	"github.com/willie68/go-micro/internal/services/adrsvc/adrmysql"
)

// New create a new storage servcice based on the configuration
func New(cfn adrsvc.Config) (adrsvc.AddressStorage, error) {
	var adrstg adrsvc.AddressStorage
	var err error
	switch cfn.Type {
	case "internal":
		adrstg, err = adrint.NewAdrInt()
		if err != nil {
			return nil, err
		}
	case "mysql":
		c := adrmysql.Config{
			Host:     cfn.Connection["host"].(string),
			Database: cfn.Connection["database"].(string),
			Table:    cfn.Connection["table"].(string),
			Username: cfn.Connection["username"].(string),
			Password: cfn.Connection["password"].(string),
		}
		adrstg, err = adrmysql.NewAdrMdb(c)
		if err != nil {
			return nil, err
		}
	}
	if adrstg == nil {
		return nil, adrsvc.ErrNotFound
	}
	do.ProvideValue(nil, adrstg)
	return adrstg, nil
}
