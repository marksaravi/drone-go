package mpu

import (
	"github.com/MarkSaravi/drone-go/types/sensore"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	InitDevice() error
	Start()
	GetDeviceConfig() (
		device sensore.Config,
		acc sensore.Config,
		gyro sensore.Config,
		err error)
	ReadRawData() ([]byte, error)
	ReadData() (acc sensore.Data, gyro sensore.Data, err error)
	WhoAmI() (string, byte, error)
	GetAcc() sensore.ThreeAxisSensore
	GetGyro() sensore.ThreeAxisSensore
}
