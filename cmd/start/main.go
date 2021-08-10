package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/motors"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	appConfig := utils.ReadConfigs()
	imuMems, esc, radio, powerBreaker := hardware.InitDroneHardware(appConfig)
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	imu := imu.CreateIM(imuMems, appConfig.Flight.Imu)
	motorsController := motors.NewMotorsControl(esc, powerBreaker)
	flightControl := flightcontrol.CreateFlightControl(imu, motorsController, radio, udpLogger)

	flightControl.Start()
}
