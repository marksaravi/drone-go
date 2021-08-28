package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/config/flightcontrolconfig"
	"github.com/MarkSaravi/drone-go/flightcontrol"
)

func main() {
	// appConfig := config.ReadConfigs()
	// imuDev, esc, radio, powerBreaker := hardware.InitDroneHardware(appConfig)
	// udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	// imu := imu.CreateIM(imuDev, appConfig.Flight.Imu)
	// motorsController := motors.NewMotorsControl(esc, powerBreaker)
	flightControlConfig := flightcontrolconfig.ReadFlightControlConfig()
	fmt.Println(flightControlConfig)
	flightControl := flightcontrol.NewFlightControl()

	flightControl.Start()
}
