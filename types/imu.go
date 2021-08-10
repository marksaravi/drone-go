package types

import "github.com/MarkSaravi/drone-go/hardware/icm20948"

type RotationsChanges struct {
	DRoll, DPitch, DYaw float64
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
		icm20948.SensorData,
		icm20948.SensorData,
		icm20948.SensorData,
		error)
}

// IMU is interface to imu mems
type IMU interface {
	GetRotations() (ImuRotations, error)
	ResetReadingTimes()
	CanRead() bool
}
