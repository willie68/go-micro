package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var prop map[string]any

func init() {
	prop = make(map[string]any)
	prop["string"] = "string value"
	prop["bool"] = true
	prop["number"] = 12345678
}

func TestConfigValueAsString(t *testing.T) {
	ast := assert.New(t)
	v, err := GetConfigValueAsString(prop, "string")
	ast.Nil(err)
	ast.Equal("string value", v)

	_, err = GetConfigValueAsString(prop, "muck")
	ast.NotNil(err)

	_, err = GetConfigValueAsString(prop, "number")
	ast.NotNil(err)

	_, err = GetConfigValueAsString(prop, "bool")
	ast.NotNil(err)
}

func TestConfigValueAsBool(t *testing.T) {
	ast := assert.New(t)
	v, err := GetConfigValueAsBool(prop, "bool")
	ast.Nil(err)
	ast.True(v)

	_, err = GetConfigValueAsBool(prop, "muck")
	ast.NotNil(err)

	_, err = GetConfigValueAsBool(prop, "number")
	ast.NotNil(err)

	_, err = GetConfigValueAsBool(prop, "string")
	ast.NotNil(err)
}

func TestConfigValueAsInt(t *testing.T) {
	ast := assert.New(t)
	v, err := GetConfigValueAsInt(prop, "number")
	ast.Nil(err)
	ast.Equal(int64(12345678), v)

	_, err = GetConfigValueAsInt(prop, "muck")
	ast.NotNil(err)

	_, err = GetConfigValueAsInt(prop, "bool")
	ast.NotNil(err)

	_, err = GetConfigValueAsInt(prop, "string")
	ast.NotNil(err)
}
