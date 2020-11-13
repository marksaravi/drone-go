package threeaxissensore

// Config is the configurations for ThreeAxisSensore
type Config interface {
}

// Data is X, Y, Z data
type Data struct {
	X, Y, Z float64
}

// ThreeAxisSensore is interface to a 3 Axis sensore
type ThreeAxisSensore interface {
	GetConfig() Config
	SetConfig(config Config)
	GetData() Data
	SetData(x, y, z float64)
	GetDiff() float64
}
