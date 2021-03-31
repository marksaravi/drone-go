package types

// Config is the configurations for ThreeAxisSensore
type Config interface {
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}
