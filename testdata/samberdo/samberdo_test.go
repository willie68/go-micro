package samberdo

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/testdata/samberdo/consumer"
	"github.com/willie68/go-micro/testdata/samberdo/producer"
)

func TestSamberDOInterfaces(t *testing.T) {
	inj := do.New()

	ast := assert.New(t)
	prod := producer.New("prodname")
	do.ProvideValue(inj, prod)

	cons := consumer.New(inj, "con")

	s := cons.Else()
	ast.NotEmpty(s)
}
