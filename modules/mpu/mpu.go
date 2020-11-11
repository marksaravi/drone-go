package mpu

import (
	"github.com/MarkSaravi/drone-go/modules/mpu/accelerometer"
	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	SetDeviceConfig() error
	GetDeviceConfig() ([]byte, error)
	ReadRawData() ([]byte, error)
	ReadData() (accX, accY, accZ, gyroX, gyroY, gyroZ float64, err error)
	WhoAmI() (string, byte, error)
	gyroscope.Gyroscope
	accelerometer.Accelerometer
}
