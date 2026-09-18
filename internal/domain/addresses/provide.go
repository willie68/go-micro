package addresses

import "github.com/samber/do/v2"

// Provide provides the Addresses instance
func Provide(inj do.Injector) error {
	do.ProvideValue(inj, NewAddresses(do.MustInvokeAs[AddressStorage](inj)))
	return nil
}
