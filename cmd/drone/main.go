package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/icm20948"
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/hardware"
)

func main() {
	// appConfig := config.ReadConfigs()
	// imuDev, esc, radio, powerBreaker := hardware.InitDroneHardware(appConfig)
	// udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	// imu := imu.CreateIM(imuDev, appConfig.Flight.Imu)
	// motorsController := motors.NewMotorsControl(esc, powerBreaker)

	config := config.ReadFlightControlConfig()
	fmt.Println(config)
	hardware.InitHost()
	imuSPIConn := drivers.NewSPIConnection(
		config.Configs.Drivers.ImuMemes.SPI.BusNumber,
		config.Configs.Drivers.ImuMemes.SPI.ChipSelect,
	)
	imuMems := icm20948.NewICM20948Driver(
		imuSPIConn,
		config.Configs.Drivers.ImuMemes.Accelerometer.SensitivityLevel,
		config.Configs.Drivers.ImuMemes.Accelerometer.Averaging,
		config.Configs.Drivers.ImuMemes.Accelerometer.LowPassFilterEnabled,
		config.Configs.Drivers.ImuMemes.Accelerometer.LowPassFilterConfig,
		config.Configs.Drivers.ImuMemes.Accelerometer.Offsets.X,
		config.Configs.Drivers.ImuMemes.Accelerometer.Offsets.Y,
		config.Configs.Drivers.ImuMemes.Accelerometer.Offsets.Z,
		config.Configs.Drivers.ImuMemes.Gyroscope.SensitivityLevel,
		config.Configs.Drivers.ImuMemes.Gyroscope.Averaging,
		config.Configs.Drivers.ImuMemes.Gyroscope.LowPassFilterEnabled,
		config.Configs.Drivers.ImuMemes.Gyroscope.LowPassFilterConfig,
		config.Configs.Drivers.ImuMemes.Gyroscope.Offsets.X,
		config.Configs.Drivers.ImuMemes.Gyroscope.Offsets.Y,
		config.Configs.Drivers.ImuMemes.Gyroscope.Offsets.Z,
	)
	readingInterval := time.Second / time.Duration(config.Configs.Devices.ImuConfig.ImuDataPerSecond)
	imu := devices.NewIMU(
		imuMems,
		readingInterval,
		config.Configs.Devices.ImuConfig.AccLowPassFilterCoefficient,
		config.Configs.Devices.ImuConfig.LowPassFilterCoefficient,
	)
	flightControl := flightcontrol.NewFlightControl(imu)

	flightControl.Start()
}
