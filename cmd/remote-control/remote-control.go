package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/mcp3008"
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
	)
	xAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Roll.Channel,
		config.RemoteControlConfig.Joysticks.Roll.ZeroValue,
	)
	yAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Pitch.Channel,
		config.RemoteControlConfig.Joysticks.Pitch.ZeroValue,
	)
	zAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfig.Joysticks.VRef,
		config.RemoteControlConfig.Joysticks.Yaw.Channel,
		config.RemoteControlConfig.Joysticks.Yaw.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisanalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisanalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisanalogToDigitalConvertor)
	gpioinput := drivers.NewGPIOSwitch(config.RemoteControlConfig.Buttons.FrontLeft)
	input := devices.NewButton(gpioinput)
	remoteControl := remotecontrol.NewRemoteControl(roll, pitch, yaw, input)
	remoteControl.Start()
}
