package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/apps/remotecontrol"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/remoteinputs"
)

func main() {
	fmt.Println("Starting RemoteControl")
	hardware.InitHost()
	btnFrontLeft := hardware.NewButton("GPIO27")
	inputs := remoteinputs.NewRemoteInputs(btnFrontLeft)
	remoteControl := remotecontrol.NewRemoteControl(inputs)
	remoteControl.Start()
}
