package config

/*
Secret our service configuration
*/
type Secret struct {
	MongoDB struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mongodb"`
}
