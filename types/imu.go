package types

// x, y, z data
type XYZ struct {
	X, Y, Z float64
}

// dx/dt, dy/dt, dz/dt
type DXYZ struct {
	DX, DY, DZ float64
}

// Inertial Measurment Unit Data (6 Degree of Freedom, Micro-electromechanical Systems)
type IMUMems6DOFRawData struct {
	Accelerometer XYZ
	Gyroscope     DXYZ
}
