package consumer

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/willie68/go-micro/testdata/samberdo/consumer/mocks"
)

func TestConsumer(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()

	m := mocks.NewInter(t)
	m.EXPECT().Doit(mock.Anything).Return("cons")

	do.ProvideValue(inj, m)

	con := New(inj, "cons")

	s := con.Else()

	ast.Equal("cons", s)
}
