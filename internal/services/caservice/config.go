package caservice

// Config the micro-vault ca service config
type Config struct {
	Servicename string `yaml:"servicename"`
	UseCA       bool   `yaml:"useca"`
	URL         string `yaml:"url"`
	AccessKey   string `yaml:"accesskey"`
	Secret      string `yaml:"secret"`
}
