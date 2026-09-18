package address

// Config general config for the storage
type Config struct {
	Type       string         `yaml:"type"`
	Connection map[string]any `yaml:"connection"`
}
