package samberdo

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/testdata/samberdo/common"
	"github.com/willie68/go-micro/testdata/samberdo/consumer"
	"github.com/willie68/go-micro/testdata/samberdo/producer1"
	"github.com/willie68/go-micro/testdata/samberdo/producer2"
)

func TestSamberDOInterfaces(t *testing.T) {
	inj := do.New()

	do.ProvideValue(inj, common.New("com"))

	inj1 := inj.Scope("producer1")
	inj2 := inj.Scope("producer2")

	ast := assert.New(t)
	prod1 := producer1.New("prodname")
	do.ProvideValue(inj1, prod1)

	prod2 := producer2.New("prodname")
	do.ProvideValue(inj2, prod2)

	cons1 := consumer.New(inj1, "con1")
	cons2 := consumer.New(inj2, "con2")

	s := cons1.Else()
	ast.NotEmpty(s)
	t.Logf("output 1: %s", s)

	s = cons2.Else()
	ast.NotEmpty(s)
	t.Logf("output 2: %s", s)
	t.Fail()
}
