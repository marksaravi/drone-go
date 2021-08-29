package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/mcp3008"
	"github.com/MarkSaravi/drone-go/remotecontrol"
)

func main() {
	fmt.Println("Starting RemoteControl")
	config := config.ReadRemoteControlConfig()
	fmt.Println(config)
	drivers.InitHost()

	analogToDigitalSPIConn := drivers.NewSPIConnection(
		config.RemoteControlConfigs.Joysticks.SPI.BusNumber,
		config.RemoteControlConfigs.Joysticks.SPI.ChipSelect,
	)
	xAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Roll.Channel,
		config.RemoteControlConfigs.Joysticks.Roll.ZeroValue,
	)
	yAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Pitch.Channel,
		config.RemoteControlConfigs.Joysticks.Pitch.ZeroValue,
	)
	zAxisanalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Yaw.Channel,
		config.RemoteControlConfigs.Joysticks.Yaw.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisanalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisanalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisanalogToDigitalConvertor)
	gpioinput := drivers.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.FrontLeft)
	input := devices.NewButton(gpioinput)
	remoteControl := remotecontrol.NewRemoteControl(roll, pitch, yaw, input)
	remoteControl.Start()
}
