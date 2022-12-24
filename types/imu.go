package types

// Raw x, y, z data
type XYZ struct {
	X, Y, Z float64
}

// Orientations (Roll, Pitch, Yaw)
type Orientations struct {
	Roll, Pitch, Yaw float64
}

// Accelerometer data (x, y, z)
type AccelerometerData struct {
	X, Y, Z float64
}

// Gyroscope data (dx/dt, dy/dt, dz/dt)
type GyroscopeData struct {
	Dx, Dy, Dz float64
}

// Inertial Measurment Unit Data (6 Degree of Freedom, Micro-electromechanical Systems)
type IMUMems6DOFData struct {
	Accelerometer AccelerometerData
	Gyroscope     GyroscopeData
}
