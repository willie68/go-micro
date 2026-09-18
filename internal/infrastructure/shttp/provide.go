package shttp

import "github.com/samber/do/v2"

func Provide(inj do.Injector) error {
	cfg := do.MustInvoke[Config](inj)
	sh, err := NewSHttp(cfg)
	if err != nil {
		return err
	}
	do.ProvideValue(inj, sh)
	return nil
}
