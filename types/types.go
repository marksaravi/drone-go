package types

// Config is the generic configuration
type Config interface {
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

// Sensor is devices that read data in x, y, z format
type Sensor struct {
	IsDataReady bool
	Data        XYZ
	Config      Config
}

// GetConfig reads the config
func (a *Sensor) GetConfig() Config {
	return a.Config
}

// SetConfig sets the config
func (a *Sensor) SetConfig(config Config) {
	a.Config = config
}

// SetData sets the data
func (a *Sensor) SetData(x, y, z float64) {
	a.IsDataReady = true
	a.Data = XYZ{
		X: x,
		Y: y,
		Z: z,
	}
}

// GetData gets the data
func (a *Sensor) GetData() (xyz XYZ, isDataReady bool) {
	return a.Data, a.IsDataReady
}
