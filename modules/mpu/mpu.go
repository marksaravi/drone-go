package mpu

import "github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	SetDeviceConfig() error
	GetDeviceConfig() ([]byte, error)
	WhoAmI() (string, byte, error)
	gyroscope.Gyroscope
}
