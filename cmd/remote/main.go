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

	var topl bool = false
	var frontl bool = false
	var downl bool = false
	var topr bool = false
	var frontr bool = false
	var downr bool = false
	for {
		rd := remoteControl.ReadInputs()
		if topr != rd.TopRight || frontr != rd.FrontRight || downr != rd.DownRight || topl != rd.TopLeft || frontl != rd.FrontLeft || downl != rd.DownLeft {
			// fmt.Printf("X: %4.1f, Y: %4.1f, Z: %4.1f, T: %4.1f, Front: %v, Top: %v, Down: %v\n", rd.X, rd.Y, rd.Z, rd.Throttle, rd.Front, rd.Top, rd.Down)
			fmt.Printf("FrontL: %v, TopL: %v, DownL: %v, FrontR: %v, TopR: %v, DownR: %v\n", rd.FrontLeft, rd.TopLeft, rd.DownLeft, rd.FrontRight, rd.TopRight, rd.DownRight)
			topl = rd.TopLeft
			downl = rd.DownLeft
			frontl = rd.FrontLeft
			topr = rd.TopRight
			downr = rd.DownRight
			frontr = rd.FrontRight
		}
		time.Sleep(time.Millisecond * 10)
	}
}
