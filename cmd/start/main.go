package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	appConfig := utils.ReadConfigs()
	imuMems, _, _ := hardware.InitHardware(appConfig)
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	imu := imu.CreateIM(imuMems, appConfig)
	pid := flightcontrol.CreatePidController()
	flightControl := flightcontrol.CreateFlightControl(imu, pid, udpLogger)

	flightControl.Start()
}
