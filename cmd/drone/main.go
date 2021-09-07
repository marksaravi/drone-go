package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {

	imu, imuDataPerSecond, escUpdatePerSecond := utils.NewImu()
	flightControl := flightcontrol.NewFlightControl(
		imuDataPerSecond,
		escUpdatePerSecond,
		imu,
		utils.NewESC(),
		utils.NewRadio(),
		utils.NewLogger(),
	)

	flightControl.Start()
}
