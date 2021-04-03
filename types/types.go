package types

// Config is the generic configuration
type Config interface {
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
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

// GetConfig reads the config
func (a *Sensor) GetConfig() Config {
	return a.Config
}

// SetConfig sets the config
func (a *Sensor) SetConfig(config Config) {
	a.Config = config
}
