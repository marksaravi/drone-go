package imu

import "github.com/MarkSaravi/drone-go/types"

const BUFFER_LEN uint8 = 4

// IMU is interface to imu mems
type IMU interface {
	Close()
	InitDevice() error
	ReadRawData() ([]byte, error)
	ReadData() (types.ImuSensorsData, error)
	GetRotations()
	WhoAmI() (string, byte, error)
}
