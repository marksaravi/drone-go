package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	devcices "github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/remotecontrol"
)

func main() {
	fmt.Println("Starting RemoteControl")
	config := config.ReadRemoteConfig()
	fmt.Println(config)
	hardware.InitHost()

	analogToDigitalSPIConn := drivers.NewSPIConnection(
		config.RemoteConfig.Joysticks.SPI.BusNumber,
		config.RemoteConfig.Joysticks.SPI.ChipSelect,
		config.RemoteConfig.Joysticks.SPI.Speed,
		config.RemoteConfig.Joysticks.SPI.Mode,
	)
	xAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteConfig.Joysticks.VRef,
		config.RemoteConfig.Joysticks.Roll.Channel,
		config.RemoteConfig.Joysticks.Roll.ZeroValue,
	)
	yAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteConfig.Joysticks.VRef,
		config.RemoteConfig.Joysticks.Pitch.Channel,
		config.RemoteConfig.Joysticks.Pitch.ZeroValue,
	)
	zAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteConfig.Joysticks.VRef,
		config.RemoteConfig.Joysticks.Yaw.Channel,
		config.RemoteConfig.Joysticks.Yaw.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisanalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisanalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisanalogToDigitalConvertor)
	gpioinput := drivers.NewGPIOSwitch(config.RemoteConfig.Buttons.FrontLeft)
	input := devcices.NewButton(gpioinput)
	remoteControl := remotecontrol.NewRemoteControl(roll, pitch, yaw, input)
	remoteControl.Start()
}
