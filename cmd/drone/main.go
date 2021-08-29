package main

import (
	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {

	flightControl := flightcontrol.NewFlightControl(
		utils.NewImu(),
		utils.NewLogger(),
	)

	flightControl.Start()
}
