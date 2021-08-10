package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/modules/adcconverter"
	"github.com/MarkSaravi/drone-go/remotecontrol"
)

func main() {
	fmt.Println("Starting Remote Control")
	config := config.ReadConfigs()
	adcDev, _ := hardware.InitRemoteHardware(config)
	adcConverter := adcconverter.NewADCConverter(adcDev)
	remoteControl := remotecontrol.NewRemoteControl(adcConverter)

	for {
		remoteControl.ReadInputs()
		time.Sleep(time.Millisecond * 1000)
	}
}
