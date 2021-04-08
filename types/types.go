package types

import "math"

// Config is the generic configuration
type Config interface {
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

func (r Rotations) Scaler() float64 {
	return math.Sqrt(r.Roll*r.Roll + r.Pitch*r.Pitch + r.Yaw*r.Yaw)
}

func toDeg(x float64) float64 {
	return x / math.Pi * 180
}

func (r *Rotations) ToDeg() Rotations {
	return Rotations{
		Roll:  toDeg(r.Roll),
		Pitch: toDeg(r.Pitch),
		Yaw:   toDeg(r.Yaw),
	}
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
