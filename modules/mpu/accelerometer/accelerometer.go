package accelerometer

// Accelerometer is interface to Accelerometer methods
type Accelerometer interface {
	GetAccelerometerConfig() ([]byte, error)
	SetAccelerometerConfig(accelerometerSensitivity int) error
}
