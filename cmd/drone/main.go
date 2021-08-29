package main

import (
	"github.com/MarkSaravi/drone-go/cmd/utils"
	"github.com/MarkSaravi/drone-go/flightcontrol"
)

func main() {

	flightControl := flightcontrol.NewFlightControl(
		utils.NewImu(),
		utils.NewLogger(),
	)

	flightControl.Start()
}
