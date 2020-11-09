package gyroscope

// Config is the configuration for Gyroscope
type Config struct {
	Scale byte
}

// Gyroscope is an interface to gyroscope mems operations
type Gyroscope interface {
	SetGyroConfig(*Config) error
	GetGyroConfig() (*Config, error)
}
