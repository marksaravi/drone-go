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
	config := config.ReadRemoteControlConfig()
	fmt.Println(config)
	hardware.InitHost()

	analogToDigitalSPIConn := drivers.NewSPIConnection(
		config.RemoteControlConfig.Joysticks.SPI.BusNumber,
		config.RemoteControlConfig.Joysticks.SPI.ChipSelect,
		config.RemoteControlConfig.Joysticks.SPI.Speed,
		config.RemoteControlConfig.Joysticks.SPI.Mode,
	)
	xAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Roll.Channel,
		config.RemoteControlConfig.Joysticks.Roll.ZeroValue,
	)
	yAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Pitch.Channel,
		config.RemoteControlConfig.Joysticks.Pitch.ZeroValue,
	)
	zAxisanalogToDigitalConvertor := drivers.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Yaw.Channel,
		config.RemoteControlConfig.Joysticks.Yaw.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisanalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisanalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisanalogToDigitalConvertor)
	gpioinput := drivers.NewGPIOSwitch(config.RemoteControlConfig.Buttons.FrontLeft)
	input := devcices.NewButton(gpioinput)
	remoteControl := remotecontrol.NewRemoteControl(roll, pitch, yaw, input)
	remoteControl.Start()
}
