package sensore

// Config is the configurations for ThreeAxisSensore
type Config interface {
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

// ThreeAxisSensore is interface to a 3 Axis sensore
type ThreeAxisSensore interface {
	GetConfig() Config
	SetConfig(config Config)
	GetData() XYZ
	SetData(x, y, z float64)
	GetDiff() float64
}
