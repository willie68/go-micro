package interfaces

// Service is the standard service interface
type Service interface {
	Init() error
	Shutdown() error
}
