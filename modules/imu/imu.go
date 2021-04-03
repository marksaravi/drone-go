package imu

import (
	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_LEN uint8 = 4

type ImuData struct {
	Acc, Gyro, Mag types.SensorData
}

// IMU is interface to imu mems
type IMU interface {
	Close() error
	InitDevice() error
	Start()
	ReadRawData() ([]byte, error)
	ReadData() (ImuData, error)
	WhoAmI() (string, byte, error)
}
