package common

import "errors"

// Error definitions
var (
	ErrNotFound = errors.New("not found")
)

// Config general config for the storage
type Config struct {
	Type       string         `yaml:"type"`
	Connection map[string]any `yaml:"connection"`
}
