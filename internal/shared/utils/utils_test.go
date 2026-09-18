package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	assert.NotEmpty(t, id)
	assert.NotEqual(t, GenerateID(), id)
}
