package health

import (
	"errors"

	"github.com/samber/do/v2"
)

// Check this is the interface for the healthcheck system. All service, which like to participate to this, should implement this interface and register via health.Register() method.
type Check interface {
	// CheckName should return the name of this healthcheck. The name should be unique.
	CheckName() string
	// Check proceed a check and return state, true for healthy or false and an optional error, if the healthcheck fails
	Check() (bool, error)
}

type Registerer interface {
	Register(check Check)
}

type Unregisterer interface {
	Unregister(checkname string) bool
}

// Register register a new healthcheck. If a healthcheck with the same name is already present, this will be overwritten
// Otherwise the new healthcheck will be appended
func Register(inj do.Injector, check Check) error {
	sh := do.MustInvokeAs[Registerer](inj)
	if sh == nil {
		return errors.New("can't get the health system service, not correctly initialised?")
	}
	sh.Register(check)
	return nil
}

// Unregister unregister a healthcheck. Return true if the healthcheck can be unregistered otherwise false
func Unregister(inj do.Injector, checkname string) error {
	sh := do.MustInvokeAs[Unregisterer](inj)
	if sh == nil {
		return errors.New("can't get the health system service, not correctly initialised?")
	}
	if !sh.Unregister(checkname) {
		return errors.New("check with name not found")
	}
	return nil
}
