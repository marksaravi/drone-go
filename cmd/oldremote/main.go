package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Remote Control")
	// config := config.ReadConfigs()
	// adcConverter, _, buttonFrontLeft, buttonFrontRight, buttonTopLeft, buttonTopRight, buttonDownLeft, buttonDownRight := hardware.InitRemoteHardware(config)
	// remoteControl := oldremotecontrol.NewRemoteControl(adcConverter, buttonFrontLeft, buttonFrontRight, buttonTopLeft, buttonTopRight, buttonDownLeft, buttonDownRight, config.RemoteControl)

	// var topl bool = false
	// var frontl bool = false
	// var downl bool = false
	// var topr bool = false
	// var frontr bool = false
	// var downr bool = false
	// for {
	// 	rd := remoteControl.ReadInputs()
	// 	if topr != rd.ButtonTopRight || frontr != rd.ButtonFrontRight || downr != rd.ButtonDownRight || topl != rd.ButtonTopLeft || frontl != rd.ButtonFrontLeft || downl != rd.ButtonDownLeft {
	// 		// fmt.Printf("X: %4.1f, Y: %4.1f, Z: %4.1f, T: %4.1f, Front: %v, Top: %v, Down: %v\n", rd.X, rd.Y, rd.Z, rd.Throttle, rd.Front, rd.Top, rd.Down)
	// 		fmt.Printf("FL: %v, FR: %v, TL: %v, TR: %v, DL: %v, DR: %v\n", rd.ButtonFrontLeft, rd.ButtonFrontRight, rd.ButtonTopLeft, rd.ButtonTopRight, rd.ButtonDownLeft, rd.ButtonDownRight)
	// 		topl = rd.ButtonTopLeft
	// 		downl = rd.ButtonDownLeft
	// 		frontl = rd.ButtonFrontLeft
	// 		topr = rd.ButtonTopRight
	// 		downr = rd.ButtonDownRight
	// 		frontr = rd.ButtonFrontRight
	// 	}
	// 	time.Sleep(time.Millisecond * 10)
	// }
}
