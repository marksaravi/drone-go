package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/apps/remotecontrol"
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
)

func main() {
	fmt.Println("Starting Remote Control")
	config := config.ReadConfigs()
	adcConverter, _, buttonFrontLeft := hardware.InitRemoteHardware(config)
	remoteControl := remotecontrol.NewRemoteControl(adcConverter, buttonFrontLeft, config.RemoteControl)

	for {
		rd := remoteControl.ReadInputs()
		fmt.Printf("X: %4.1f, Y: %4.1f, Z: %4.1f, T: %4.1f, Stop: %v\n", rd.X, rd.Y, rd.Z, rd.Throttle, rd.Stop)
		time.Sleep(time.Millisecond * 10)
	}
}
