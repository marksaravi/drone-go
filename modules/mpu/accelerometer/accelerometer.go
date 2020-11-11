package accelerometer

// Config is the configurations for Accelerometer
type Config struct {
	Sensitivity int
}

// Accelerometer is interface to Accelerometer methods
type Accelerometer interface {
	GetAccelerometerConfig() ([]byte, error)
	SetAccelerometerConfig(config Config) error
}
