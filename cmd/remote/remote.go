package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/apps/remotecontrol"
	"github.com/MarkSaravi/drone-go/hardware"
)

func main() {
	fmt.Println("Starting RemoteControl")
	hardware.InitHost()
	remoteControl := remotecontrol.NewRemoteControl(remotecontrol.RemoteConfig{
		ButtonFrontEndGPIO: "GPIO27",
	})
	remoteControl.Start()
}
