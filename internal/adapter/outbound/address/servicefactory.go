package address

import (
	"github.com/samber/do/v2"
	"github.com/willie68/go-micro/internal/adapter/outbound/address/common"
	adrint "github.com/willie68/go-micro/internal/adapter/outbound/address/memory"
	adrmysql "github.com/willie68/go-micro/internal/adapter/outbound/address/mysql"
)

// Provide create a new storage servcice based on the configuration
func Provide(inj do.Injector) error {
	cfg := do.MustInvoke[Config](inj)
	switch cfg.Type {
	case "internal":
		adrstg, err := adrint.NewAdrInt()
		if err != nil {
			return err
		}
		do.ProvideValue(inj, adrstg)
		return err
	case "mysql":
		c := adrmysql.Config{
			Host:     cfg.Connection["host"].(string),
			Database: cfg.Connection["database"].(string),
			Table:    cfg.Connection["table"].(string),
			Username: cfg.Connection["username"].(string),
			Password: cfg.Connection["password"].(string),
		}
		sqlstg, err := adrmysql.NewAdrMdb(c)
		if err != nil {
			return err
		}
		do.ProvideValue(inj, sqlstg)
		return err
	}
	return common.ErrNotFound
}
