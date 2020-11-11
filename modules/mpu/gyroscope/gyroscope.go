package gyroscope

// Config is the configuration for Gyroscope
type Config struct {
	FullScale int
}

// Gyroscope is an interface to gyroscope mems operations
type Gyroscope interface {
	SetGyroConfig(*Config) error
	GetGyroConfig() (*Config, error)
}
