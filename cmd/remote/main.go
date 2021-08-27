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
	adcConverter, _, buttonFrontLeft, buttonFrontRight, buttonTopLeft, buttonTopRight, buttonDownLeft, buttonDownRight := hardware.InitRemoteHardware(config)
	remoteControl := remotecontrol.NewRemoteControl(adcConverter, buttonFrontLeft, buttonFrontRight, buttonTopLeft, buttonTopRight, buttonDownLeft, buttonDownRight, config.RemoteControl)

	var top bool = false
	var front bool = false
	var down bool = false
	for {
		rd := remoteControl.ReadInputs()
		if top != rd.Top || front != rd.Front || down != rd.Down {
			fmt.Printf("X: %4.1f, Y: %4.1f, Z: %4.1f, T: %4.1f, Front: %v, Top: %v, Down: %v\n", rd.X, rd.Y, rd.Z, rd.Throttle, rd.Front, rd.Top, rd.Down)
			top = rd.Top
			down = rd.Down
			front = rd.Front
		}
		time.Sleep(time.Millisecond * 10)
	}
}
