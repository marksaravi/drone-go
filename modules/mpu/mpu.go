package mpu

import (
	"github.com/MarkSaravi/drone-go/modules/mpu/gyroscope"
)

// MPU is interface to mpu mems
type MPU interface {
	Close()
	ResetToDefault()
	gyroscope.Gyroscope
}
