package mpu

import (
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/types/sensore"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	InitDevice() error
	Start()
	GetDeviceConfig() (
		device types.Config,
		acc types.Config,
		gyro types.Config,
		err error)
	ReadRawData() ([]byte, error)
	ReadData() (acc types.XYZ, gyro types.XYZ, err error)
	WhoAmI() (string, byte, error)
	GetAcc() sensore.ThreeAxisSensore
	GetGyro() sensore.ThreeAxisSensore
}
