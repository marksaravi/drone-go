package mpu

import (
	"github.com/MarkSaravi/drone-go/types"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	InitDevice() error
	Start()
	ReadRawData() ([]byte, error)
	ReadData() (acc types.XYZ, isAccDataReady bool, gyro types.XYZ, isGyroDataReady bool, err error)
	WhoAmI() (string, byte, error)
}
