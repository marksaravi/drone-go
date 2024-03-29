package mems

import "time"

// XYZ: x, y, z data
type XYZ struct {
	X, Y, Z float64
}

// DXYZ: dx/dt, dy/dt, dz/dt
type DXYZ struct {
	DX, DY, DZ float64
}

// Inertial Measurment Unit Data (6 Degree of Freedom, Micro-electromechanical Systems)
type Mems6DOFData struct {
	Accelerometer XYZ
	Gyroscope     DXYZ
	Time          time.Time
}
