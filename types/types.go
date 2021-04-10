package types

// Config is the generic configuration
type Config interface {
}

type FlightConfig struct {
	AccLowPassFilterCoefficient       float64
	RotationsLowPassFilterCoefficient float64
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type SensorData struct {
	Error error
	Data  XYZ
}

// Sensor is devices that read data in x, y, z format
type Sensor struct {
	Type   string
	Config Config
}

// CommandParameters is parameters for the command
type CommandParameters interface {
}

type Command struct {
	Command    string
	Parameters CommandParameters
}

// GetConfig reads the config
func (a *Sensor) GetConfig() Config {
	return a.Config
}

// SetConfig sets the config
func (a *Sensor) SetConfig(config Config) {
	a.Config = config
}

type Offsets struct {
	X float64
	Y float64
	Z float64
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	SensitivityLevel int
	Offsets          []Offsets
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	ScaleLevel             int
	LowPassFilterEnabled   bool
	LowPassFilter          int
	LowPassFilterAveraging int
	Offsets                []Offsets
}

// MagnetometerConfig is the configuration for Magnetometer
type MagnetometerConfig struct {
}

type ApplicationConfig struct {
	Acc string `yaml:"host"`
}
