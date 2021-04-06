package imu

import (
	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_LEN uint8 = 4

type ImuData struct {
	Duration       float64
	Acc, Gyro, Mag types.SensorData
}

// IMU is interface to imu mems
type IMU interface {
	Close()
	InitDevice() error
	ResetGyroTimer()
	ReadRawData() ([]byte, error)
	ReadData() (ImuData, error)
	WhoAmI() (string, byte, error)
}
