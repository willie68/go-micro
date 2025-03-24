package consumer

import (
	"fmt"

	"github.com/samber/do/v2"
	"github.com/willie68/go-micro/testdata/samberdo/common"
)

//go:generate mockery --name=Inter --outpkg=mocks --filename=inter.go --with-expecter
type Inter interface {
	Doit(in string) string
}

type Consumer struct {
	name  string
	inter Inter
	com   *common.Common
}

func New(inj do.Injector, name string) *Consumer {
	return &Consumer{
		name:  name,
		inter: do.MustInvokeAs[Inter](inj),
		com:   do.MustInvoke[*common.Common](inj),
	}
}

func (c *Consumer) Else() string {
	return fmt.Sprintf("%s %s", c.com.Name(), c.inter.Doit(c.name))
}
