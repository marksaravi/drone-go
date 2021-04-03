package main

import (
	"os"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
)

func initiateIMU() *icm20948.Device {
	dev, err := icm20948.NewICM20948Driver(icm20948.Settings{
		BusNumber:  0,
		ChipSelect: 0,
		Config:     icm20948.DeviceConfig{},
		AccConfig:  icm20948.AccelerometerConfig{SensitivityLevel: 3},
		GyroConfig: icm20948.GyroscopeConfig{ScaleLevel: 2},
	},
	)
	if err != nil {
		os.Exit(1)
	}
	if dev.InitDevice() != nil {
		os.Exit(1)
	}
	return dev
}
