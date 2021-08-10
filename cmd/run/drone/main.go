package main

import (
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/motors"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
)

func main() {
	appConfig := config.ReadConfigs()
	imuDev, esc, radio, powerBreaker := hardware.InitDroneHardware(appConfig)
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	imu := imu.CreateIM(imuDev, appConfig.Flight.Imu)
	motorsController := motors.NewMotorsControl(esc, powerBreaker)
	flightControl := flightcontrol.CreateFlightControl(imu, motorsController, radio, udpLogger)

	flightControl.Start()
}
