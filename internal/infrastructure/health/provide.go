package health

import "github.com/samber/do/v2"

func Provide(inj do.Injector) error {
	health, err := NewHealthSystem(inj, do.MustInvoke[Config](inj))
	if err != nil {
		return err
	}
	do.ProvideValue(inj, health)
	return nil
}
