package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/remotecontrol"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	log.Println("Starting RemoteControl")
	config := config.ReadRemoteControlConfig()
	hardware.InitHost()

	radio := nrf204.NewRadio()
	analogToDigitalSPIConn := hardware.NewSPIConnection(
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
	gpioFrontLeft := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.FrontLeft)
	btnFrontLeft := devices.NewButton(gpioFrontLeft)
	gpioFrontRight := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.FrontRight)
	btnFrontRight := devices.NewButton(gpioFrontRight)
	gpioTopLeft := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.TopLeft)
	btnToptLeft := devices.NewButton(gpioTopLeft)
	gpioTopRight := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.TopRight)
	btnTopRight := devices.NewButton(gpioTopRight)
	gpioBottomLeft := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.BottomLeft)
	btnBottomLeft := devices.NewButton(gpioBottomLeft)
	gpioBottomRight := hardware.NewGPIOSwitch(config.RemoteControlConfigs.Buttons.BottomRight)
	btnBottomRight := devices.NewButton(gpioBottomRight)

	ctx, cancel := context.WithCancel(context.Background())
	remoteControl := remotecontrol.NewRemoteControl(
		radio,
		roll, pitch, yaw, throttle,
		btnFrontLeft, btnFrontRight,
		btnToptLeft, btnTopRight,
		btnBottomLeft, btnBottomRight,
	)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer log.Println("Stopping the Remote Control")

		log.Println("Press ENTER to Stop")
		fmt.Scanln()
		cancel()
	}(&waitGroup)
	remoteControl.Start(ctx, &waitGroup)
	waitGroup.Wait()
	log.Println("Remote Control stopped")
}
