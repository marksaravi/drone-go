package main

import (
	"github.com/MarkSaravi/drone-go/devices/icm20948"
	"github.com/MarkSaravi/drone-go/types"
)

type ApplicationConfig struct {
	FlightConfig types.FlightConfig
	ICM20948     struct {
		BusNumber     int
		ChipSelect    int
		Accelerometer icm20948.AccelerometerConfig
		Gyroscope     icm20948.GyroscopeConfig
		Magnetometer  icm20948.MagnetometerConfig
	}
}
