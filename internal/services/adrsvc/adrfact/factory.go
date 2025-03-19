package adrfact

import (
	"github.com/samber/do/v2"
	"github.com/willie68/go-micro/internal/services/adrsvc"
	"github.com/willie68/go-micro/internal/services/adrsvc/adrint"
	"github.com/willie68/go-micro/internal/services/adrsvc/adrmysql"
	"github.com/willie68/go-micro/internal/services/health"
)

// New create a new storage servcice based on the configuration
func New(inj do.Injector, cfn adrsvc.Config) error {
	switch cfn.Type {
	case "internal":
		adrstg, err := adrint.NewAdrInt()
		if err != nil {
			return err
		}
		do.ProvideValue(inj, adrstg)

		err = health.Register(inj, adrstg)
		return err
	case "mysql":
		c := adrmysql.Config{
			Host:     cfn.Connection["host"].(string),
			Database: cfn.Connection["database"].(string),
			Table:    cfn.Connection["table"].(string),
			Username: cfn.Connection["username"].(string),
			Password: cfn.Connection["password"].(string),
		}
		sqlstg, err := adrmysql.NewAdrMdb(c)
		if err != nil {
			return err
		}
		do.ProvideValue(inj, sqlstg)

		err = health.Register(inj, sqlstg)
		return err
	}
	return adrsvc.ErrNotFound
}
