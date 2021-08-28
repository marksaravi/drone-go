package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/apps/remotecontrol"
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/hardware/mcp3008"
	"github.com/MarkSaravi/drone-go/modules/remoteinputs"
)

func main() {
	fmt.Println("Starting RemoteControl")
	config := config.ReadRemoteConfig()
	fmt.Println(config)
	hardware.InitHost()
	btnFrontLeft := hardware.NewButton(config.RemoteConfig.Buttons.FrontLeft)
	analogToDigitalConvertor := mcp3008.NewMCP3008(
		config.RemoteConfig.Joysticks.SPI.BusNumber,
		config.RemoteConfig.Joysticks.SPI.ChipSelect,
		config.RemoteConfig.Joysticks.SPI.Mode,
		config.RemoteConfig.Joysticks.SPI.Speed,
	)
	roll := hardware.NewJoystick(
		analogToDigitalConvertor,
		config.RemoteConfig.Joysticks.Roll.Channel,
		config.RemoteConfig.Joysticks.Roll.ZeroValue,
		config.RemoteConfig.Joysticks.VRef,
	)
	inputs := remoteinputs.NewRemoteInputs(roll, btnFrontLeft)
	remoteControl := remotecontrol.NewRemoteControl(inputs)
	remoteControl.Start()
}
