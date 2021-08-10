package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/remotecontrol"
)

func main() {
	fmt.Println("Starting Remote Control")
	config := config.ReadConfigs()
	adcConverter, _ := hardware.InitRemoteHardware(config)
	remoteControl := remotecontrol.NewRemoteControl(adcConverter, config.RemoteControl)

	for {
		rd := remoteControl.ReadInputs()
		fmt.Println(rd)
		time.Sleep(time.Millisecond * 1000)
	}
}
