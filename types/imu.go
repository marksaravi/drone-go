package types

// Raw x, y, z data
type XYZ struct {
	X, Y, Z float64
}

type XYZDt struct {
	DX, DY, DZ float64
}

// Orientations (Roll, Pitch, Yaw)
type Rotations struct {
	Roll, Pitch, Yaw float64
}

// Inertial Measurment Unit Data (6 Degree of Freedom, Micro-electromechanical Systems)
type IMUMems6DOFRawData struct {
	Accelerometer XYZ
	Gyroscope     XYZDt
}
