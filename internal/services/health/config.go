package health

// Config configuration for the healthcheck system
type Config struct {
	// Period in seconds, when all health services should run
	Period int `yaml:"period"`
	// StartDelay an optional starting delay, after starting the service
	StartDelay int `yaml:"startdelay"`
}
