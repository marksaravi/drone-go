package mpu

import (
	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	InitDevice() error
	Start()
	GetDeviceConfig() (
		device threeaxissensore.Config,
		acc threeaxissensore.Config,
		gyro threeaxissensore.Config,
		err error)
	ReadRawData() ([]byte, error)
	ReadData() (acc threeaxissensore.Data, gyro threeaxissensore.Data, err error)
	WhoAmI() (string, byte, error)
	GetAcc() threeaxissensore.ThreeAxisSensore
	GetGyro() threeaxissensore.ThreeAxisSensore
}
