package imu

import "github.com/MarkSaravi/drone-go/types"

const BUFFER_LEN uint8 = 4

// IMU is interface to imu mems
type IMU interface {
	Close()
	InitDevice() error
	ReadSensorsRawData() ([]byte, error)
	ReadSensors() (types.ImuSensorsData, error)
	GetRotations() (types.ImuRotations, error)
	WhoAmI() (string, byte, error)
}
