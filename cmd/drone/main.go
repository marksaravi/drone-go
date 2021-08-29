package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/devices/udplogger"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/icm20948"
	"github.com/MarkSaravi/drone-go/flightcontrol"
)

func main() {
	config := config.ReadFlightControlConfig()
	imuConfig := config.Configs.Imu
	fmt.Println(config)
	drivers.InitHost()
	imuSPIConn := drivers.NewSPIConnection(
		imuConfig.SPI.BusNumber,
		imuConfig.SPI.ChipSelect,
	)
	accConfig := imuConfig.Accelerometer
	gyroConfig := imuConfig.Gyroscope
	imuMems := icm20948.NewICM20948Driver(
		imuSPIConn,
		accConfig.SensitivityLevel,
		accConfig.Averaging,
		accConfig.LowPassFilterEnabled,
		accConfig.LowPassFilterConfig,
		accConfig.Offsets.X,
		accConfig.Offsets.Y,
		accConfig.Offsets.Z,
		gyroConfig.SensitivityLevel,
		gyroConfig.Averaging,
		gyroConfig.LowPassFilterEnabled,
		gyroConfig.LowPassFilterConfig,
		gyroConfig.Offsets.X,
		gyroConfig.Offsets.Y,
		gyroConfig.Offsets.Z,
	)
	readingInterval := time.Second / time.Duration(imuConfig.ImuDataPerSecond)
	imu := devices.NewIMU(
		imuMems,
		readingInterval,
		imuConfig.AccLowPassFilterCoefficient,
		imuConfig.LowPassFilterCoefficient,
	)
	udplogger := udplogger.NewUdpLogger(
		true,
		"192.168.1.101",
		6431,
		20,
		20,
		imuConfig.ImuDataPerSecond,
	)
	flightControl := flightcontrol.NewFlightControl(imu, udplogger)

	flightControl.Start()
}
