package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/apps/remotecontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	log.Println("Starting RemoteControl")
	configs := config.ReadRemoteControlConfig()
	fmt.Println(configs)
	remoteConfigs := configs.Configs
	radioConnectionConfigs := configs.RadioConnection
	hardware.InitHost()

	radioNRF204 := nrf204.NewRadio(remoteConfigs.Radio, radioConnectionConfigs)
	radioDev := radio.NewRadio(radioNRF204, radioConnectionConfigs.ConnectionTimeoutMS)
	analogToDigitalSPIConn := hardware.NewSPIConnection(
		remoteConfigs.Joysticks.SPI.BusNumber,
		remoteConfigs.Joysticks.SPI.ChipSelect,
	)
	xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		remoteConfigs.Joysticks.VRef,
		remoteConfigs.Joysticks.Roll.Channel,
		remoteConfigs.Joysticks.Roll.ZeroValue,
	)
	yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		remoteConfigs.Joysticks.VRef,
		remoteConfigs.Joysticks.Pitch.Channel,
		remoteConfigs.Joysticks.Pitch.ZeroValue,
	)
	zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		remoteConfigs.Joysticks.VRef,
		remoteConfigs.Joysticks.Yaw.Channel,
		remoteConfigs.Joysticks.Yaw.ZeroValue,
	)
	throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		remoteConfigs.Joysticks.VRef,
		remoteConfigs.Joysticks.Throttle.Channel,
		remoteConfigs.Joysticks.Throttle.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisAnalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisAnalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisAnalogToDigitalConvertor)
	throttle := devices.NewJoystick(throttleAlogToDigitalConvertor)
	gpioFrontLeft := hardware.NewGPIOSwitch(remoteConfigs.Buttons.FrontLeft)
	btnFrontLeft := devices.NewButton(gpioFrontLeft)
	gpioFrontRight := hardware.NewGPIOSwitch(remoteConfigs.Buttons.FrontRight)
	btnFrontRight := devices.NewButton(gpioFrontRight)
	gpioTopLeft := hardware.NewGPIOSwitch(remoteConfigs.Buttons.TopLeft)
	btnToptLeft := devices.NewButton(gpioTopLeft)
	gpioTopRight := hardware.NewGPIOSwitch(remoteConfigs.Buttons.TopRight)
	btnTopRight := devices.NewButton(gpioTopRight)
	gpioBottomLeft := hardware.NewGPIOSwitch(remoteConfigs.Buttons.BottomLeft)
	btnBottomLeft := devices.NewButton(gpioBottomLeft)
	gpioBottomRight := hardware.NewGPIOSwitch(remoteConfigs.Buttons.BottomRight)
	btnBottomRight := devices.NewButton(gpioBottomRight)

	remoteControl := remotecontrol.NewRemoteControl(
		radioDev,
		roll, pitch, yaw, throttle,
		btnFrontLeft, btnFrontRight,
		btnToptLeft, btnTopRight,
		btnBottomLeft, btnBottomRight,
		remoteConfigs.CommandPerSecond,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var waitGroup sync.WaitGroup

	radioDev.Start(ctx, &waitGroup)
	remoteControl.Start(ctx, &waitGroup)
	utils.WaitToAbortByENTER(cancel, &waitGroup)
	waitGroup.Wait()
	log.Println("Remote Control stopped")
}
