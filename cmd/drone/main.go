package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {

	flightControl := flightcontrol.NewFlightControl(
		utils.NewImu(),
		utils.NewRadio(),
		utils.NewLogger(),
	)

	flightControl.Start()
}
