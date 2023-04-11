package services

import (
	"github.com/willie68/go-micro/internal/config"
	"github.com/willie68/go-micro/internal/services/shttp"
)

// InitServices initialise the service system
func InitServices(cfg config.Config) error {
	// here you can add more services

	_, err := shttp.NewSHttp(cfg)
	if err != nil {
		return err
	}

	return nil
}
