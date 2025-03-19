package consumer

import (
	"github.com/samber/do/v2"
)

//go:generate mockery --name=Inter --outpkg=mocks --filename=inter.go --with-expecter
type Inter interface {
	Doit(in string) string
}

type Consumer struct {
	name  string
	inter Inter
}

func New(inj do.Injector, name string) *Consumer {
	return &Consumer{
		name:  name,
		inter: do.MustInvokeAs[Inter](inj),
	}
}

func (c *Consumer) Else() string {
	return c.inter.Doit(c.name)
}
