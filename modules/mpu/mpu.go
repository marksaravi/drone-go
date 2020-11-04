package mpu

import (
	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
)

// MPU is interface to mpu mems
type MPU interface {
	Close()
	WhoAmI() (byte, error)
	gyroscope.Gyroscope
}

//NewMPU creates an ampu device based on icm20948 mems
func NewMPU(spidevice string) (MPU, error) {
	driver, err := icm20948.NewRaspberryPiICM20948Driver(spidevice)
	return driver, err
}
