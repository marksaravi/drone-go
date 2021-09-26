package main

import (
	"context"
	"fmt"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devicecreators"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/drivers/mcp3008"
	"github.com/marksaravi/drone-go/remotecontrol"
)

func main() {
	fmt.Println("Starting RemoteControl")
	config := config.ReadRemoteControlConfig()
	fmt.Println(config)
	drivers.InitHost()

	radio := devicecreators.NewRadio()
	analogToDigitalSPIConn := drivers.NewSPIConnection(
		config.RemoteControlConfigs.Joysticks.SPI.BusNumber,
		config.RemoteControlConfigs.Joysticks.SPI.ChipSelect,
	)
	xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Roll.Channel,
		config.RemoteControlConfigs.Joysticks.Roll.ZeroValue,
	)
	yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Pitch.Channel,
		config.RemoteControlConfigs.Joysticks.Pitch.ZeroValue,
	)
	zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Yaw.Channel,
		config.RemoteControlConfigs.Joysticks.Yaw.ZeroValue,
	)
	throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		config.RemoteControlConfigs.Joysticks.VRef,
		config.RemoteControlConfigs.Joysticks.Throttle.Channel,
		config.RemoteControlConfigs.Joysticks.Throttle.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisAnalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisAnalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisAnalogToDigitalConvertor)
	throttle := devices.NewJoystick(throttleAlogToDigitalConvertor)
	gpioinput := drivers.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.FrontLeft)
	input := devices.NewButton(gpioinput)
	remoteControl := remotecontrol.NewRemoteControl(radio, roll, pitch, yaw, throttle, input)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("Press ENTER to Stop")
		fmt.Scanln()
		cancel()
	}()
	remoteControl.Start(ctx)
	<-ctx.Done()
}
