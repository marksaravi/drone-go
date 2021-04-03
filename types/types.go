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
