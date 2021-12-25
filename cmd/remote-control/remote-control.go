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
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	log.Println("Starting RemoteControl")
	configs := config.ReadConfigs().RemoteControl
	fmt.Println(configs)
	radioConfigs := configs.Radio
	joysticksConfigs := configs.Joysticks
	buttonsConfis := configs.Buttons
	hardware.InitHost()

	radioNRF204 := nrf204.NewRadio(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	radioDev := radio.NewRadio(radioNRF204, radioConfigs.HeartBeatTimeoutMS)
	analogToDigitalSPIConn := hardware.NewSPIConnection(
		joysticksConfigs.SPI.BusNumber,
		joysticksConfigs.SPI.ChipSelect,
	)
	xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.VRef,
		joysticksConfigs.Roll.Channel,
		joysticksConfigs.Roll.ZeroValue,
	)
	yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.VRef,
		joysticksConfigs.Pitch.Channel,
		joysticksConfigs.Pitch.ZeroValue,
	)
	zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.VRef,
		joysticksConfigs.Yaw.Channel,
		joysticksConfigs.Yaw.ZeroValue,
	)
	throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.VRef,
		joysticksConfigs.Throttle.Channel,
		joysticksConfigs.Throttle.ZeroValue,
	)
	roll := devices.NewJoystick(xAxisAnalogToDigitalConvertor)
	pitch := devices.NewJoystick(yAxisAnalogToDigitalConvertor)
	yaw := devices.NewJoystick(zAxisAnalogToDigitalConvertor)
	throttle := devices.NewJoystick(throttleAlogToDigitalConvertor)
	gpioFrontLeft := hardware.NewGPIOSwitch(buttonsConfis.FrontLeft)
	btnFrontLeft := devices.NewButton(gpioFrontLeft)
	gpioFrontRight := hardware.NewGPIOSwitch(buttonsConfis.FrontRight)
	btnFrontRight := devices.NewButton(gpioFrontRight)
	gpioTopLeft := hardware.NewGPIOSwitch(buttonsConfis.TopLeft)
	btnToptLeft := devices.NewButton(gpioTopLeft)
	gpioTopRight := hardware.NewGPIOSwitch(buttonsConfis.TopRight)
	btnTopRight := devices.NewButton(gpioTopRight)
	gpioBottomLeft := hardware.NewGPIOSwitch(buttonsConfis.BottomLeft)
	btnBottomLeft := devices.NewButton(gpioBottomLeft)
	gpioBottomRight := hardware.NewGPIOSwitch(buttonsConfis.BottomRight)
	btnBottomRight := devices.NewButton(gpioBottomRight)
	var buzzerPin gpio.PinOut = gpioreg.ByName(configs.BuzzerGPIO)
	buzzer := piezobuzzer.NewBuzzer(buzzerPin)

	remoteControl := remotecontrol.NewRemoteControl(
		radioDev,
		roll, pitch, yaw, throttle,
		btnFrontLeft, btnFrontRight,
		btnToptLeft, btnTopRight,
		btnBottomLeft, btnBottomRight,
		configs.CommandPerSecond,
		buzzer,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var waitGroup sync.WaitGroup

	radioDev.Start(ctx, &waitGroup)
	remoteControl.Start(ctx, &waitGroup)
	utils.WaitToAbortByENTER(cancel, &waitGroup)
	waitGroup.Wait()
	log.Println("Remote Control stopped")
}
