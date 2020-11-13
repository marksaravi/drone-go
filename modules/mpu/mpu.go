package mpu

import (
	"github.com/MarkSaravi/drone-go/modules/mpu/accelerometer"
	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	SetDeviceConfig() error
	Start()
	GetDeviceConfig() ([]byte, error)
	ReadRawData() ([]byte, error)
	ReadData() (acc threeaxissensore.Data, gyro threeaxissensore.Data, err error)
	WhoAmI() (string, byte, error)
	gyroscope.Gyroscope
	accelerometer.Accelerometer
	GetAcc() threeaxissensore.ThreeAxisSensore
	GetGyro() threeaxissensore.ThreeAxisSensore
}
