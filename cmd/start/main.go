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
	imuMems, esc, powerBreaker := hardware.InitHardware(appConfig)
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	imu := imu.CreateIM(imuMems, appConfig.Flight.Imu)
	motorsController := motors.NewMotorsControl(esc, powerBreaker)
	pid := flightcontrol.CreatePidController()
	flightControl := flightcontrol.CreateFlightControl(imu, motorsController, pid, udpLogger)

	flightControl.Start()
}
