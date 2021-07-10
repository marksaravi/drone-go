package types

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
