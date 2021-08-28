package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/apps/remotecontrol"
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/remoteinputs"
)

func main() {
	fmt.Println("Starting RemoteControl")
	config := config.ReadRemoteConfig()
	fmt.Println(config)
	hardware.InitHost()
	btnFrontLeft := hardware.NewButton(config.RemoteConfig.Buttons.FrontLeft)
	inputs := remoteinputs.NewRemoteInputs(btnFrontLeft)
	remoteControl := remotecontrol.NewRemoteControl(inputs)
	remoteControl.Start()
}
