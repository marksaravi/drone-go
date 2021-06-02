package imu

import (
	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_LEN uint8 = 4

type ImuData struct {
	Acc, Gyro, Mag types.SensorData
	SampleRate     int
	ReadTime       int64
	ReadInterval   int64
}

// IMU is interface to imu mems
type IMU interface {
	Close()
	InitDevice() error
	ReadRawData() ([]byte, error)
	ReadData() (ImuData, error)
	WhoAmI() (string, byte, error)
}
