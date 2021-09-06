package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {

	imu, imuDataPerSecond := utils.NewImu()
	flightControl := flightcontrol.NewFlightControl(
		imuDataPerSecond,
		imu,
		utils.NewRadio(),
		utils.NewLogger(),
	)

	flightControl.Start()
}
