package main

import (
	"os"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/modules/imu"
)

func initiateIMU(config icm20948.Config, lowPassFilterCoefficient float64) imu.IMU {
	dev, err := icm20948.NewICM20948Driver(config, lowPassFilterCoefficient)
	if err != nil {
		os.Exit(1)
	}
	if dev.InitDevice() != nil {
		os.Exit(1)
	}
	return dev
}
