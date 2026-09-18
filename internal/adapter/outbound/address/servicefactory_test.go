package address

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/adapter/outbound/address/common"
)

func TestProvideInternal(t *testing.T) {
	inj := do.New()
	do.ProvideValue(inj, Config{Type: "internal"})
	assert.NoError(t, Provide(inj))
}

func TestProvideUnknown(t *testing.T) {
	inj := do.New()
	do.ProvideValue(inj, Config{Type: "unknown"})
	assert.ErrorIs(t, Provide(inj), common.ErrNotFound)
}
