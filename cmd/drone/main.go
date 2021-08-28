package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/flightcontrol"
)

func main() {
	// appConfig := config.ReadConfigs()
	// imuDev, esc, radio, powerBreaker := hardware.InitDroneHardware(appConfig)
	// udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	// imu := imu.CreateIM(imuDev, appConfig.Flight.Imu)
	// motorsController := motors.NewMotorsControl(esc, powerBreaker)

	config := config.ReadFlightControlConfig()
	fmt.Println(config)
	readingInterval := time.Second / time.Duration(config.Configs.Devices.ImuConfig.ImuDataPerSecond)
	imu := devices.NewIMU(readingInterval)
	flightControl := flightcontrol.NewFlightControl(imu)

	flightControl.Start()
}
