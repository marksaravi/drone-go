package types

// Register is the address and bank of the Register
type Register struct {
	Address uint8
	Bank    uint8
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type RotationsChanges struct {
	DRoll, DPitch, DYaw float64
}

type SensorData struct {
	Error error
	Data  XYZ
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type ImuRotations struct {
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	ReadTime      int64
	ReadInterval  int64
}

type ImuMems interface {
	ReadSensors() (
		SensorData,
		SensorData,
		SensorData,
		error)
}

// IMU is interface to imu mems
type IMU interface {
	GetRotations() (ImuRotations, error)
	ResetReadingTimes()
	CanRead() bool
}
