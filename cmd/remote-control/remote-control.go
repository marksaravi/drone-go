package main

import (
	"context"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/apps/remote"
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting RemoteControl")
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	log.Println(configs)

	radioConfigs := configs.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI.BusNum,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.SPI.SpiChipEnabledGPIO,
		radioConfigs.RxTxAddress,
	)

	radioTransmitter := radio.NewTransmitter(radioLink)

	analogToDigitalSPIConn := hardware.NewSPIConnection(
		configs.Joysticks.SPI.BusNum,
		configs.Joysticks.SPI.ChipSelect,
	)

	joystickRoll := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.RollChannel,
		configs.Joysticks.RollMidValue,
		constants.JOYSTICK_RANGE_DEG,
	)
	joystickPitch := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.PitchChannel,
		configs.Joysticks.PitchMidValue,
		constants.JOYSTICK_RANGE_DEG,
	)
	joystickYaw := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.YawChannel,
		configs.Joysticks.YawMidValue,
		constants.JOYSTICK_RANGE_DEG,
	)
	joystickThrottle := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.ThrottleChannel,
		mcp3008.DIGITAL_MAX_VALUE/2,
		constants.THROTTLE_MAX,
	)

	remoteControl := remote.NewRemoteControl(remote.RemoteSettings{
		Transmitter:      radioTransmitter,
		CommandPerSecond: configs.CommandsPerSecond,
		Roll:             joystickRoll,
		Pitch:            joystickPitch,
		Yaw:              joystickYaw,
		Throttle:         joystickThrottle,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Scanln()
		cancel()
	}()

	remoteControl.Start(ctx)

	// oledConn, err := i2creg.Open("")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer oledConn.Close()
	// oledDev := &i2c.Dev{Addr: configs.DisplayAddress, Bus: oledConn}
	// oled := ssd1306.NewSSD1306(oledDev, ssd1306.DefaultOptions)
	// oled.Init()

	// radioNRF24L01 := nrf24l01.NewNRF24L01EnhancedBurst(
	// 	radioConfigs.SPI.BusNumber,
	// 	radioConfigs.SPI.ChipSelect,
	// 	radioConfigs.CE,
	// 	radioConfigs.RxTxAddress,
	// )
	// radioDev := radio.NewTransmitter(radioNRF24L01, radioConfigs.ConnectionTimeoutMs)
	// analogToDigitalSPIConn := hardware.NewSPIConnection(
	// 	joysticksConfigs.SPI.BusNumber,
	// 	joysticksConfigs.SPI.ChipSelect,
	// )
	// xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
	// 	analogToDigitalSPIConn,
	// 	joysticksConfigs.Roll.Channel,
	// )
	// yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
	// 	analogToDigitalSPIConn,
	// 	joysticksConfigs.Pitch.Channel,
	// )
	// zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
	// 	analogToDigitalSPIConn,
	// 	joysticksConfigs.Yaw.Channel,
	// )
	// throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
	// 	analogToDigitalSPIConn,
	// 	joysticksConfigs.Throttle.Channel,
	// )
	// roll := devices.NewJoystick(xAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Roll.Offset, joysticksConfigs.Roll.Dir)
	// pitch := devices.NewJoystick(yAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Pitch.Offset, joysticksConfigs.Pitch.Dir)
	// yaw := devices.NewJoystick(zAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Yaw.Offset, joysticksConfigs.Yaw.Dir)
	// throttle := devices.NewJoystick(throttleAlogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Throttle.Offset, 1)
	// gpioFrontLeft := hardware.NewGPIOSwitch(buttonsConfis.FrontLeft)
	// btnFrontLeft := devices.NewButton(gpioFrontLeft)
	// gpioFrontRight := hardware.NewGPIOSwitch(buttonsConfis.FrontRight)
	// btnFrontRight := devices.NewButton(gpioFrontRight)
	// gpioTopLeft := hardware.NewGPIOSwitch(buttonsConfis.TopLeft)
	// btnToptLeft := devices.NewButton(gpioTopLeft)
	// gpioTopRight := hardware.NewGPIOSwitch(buttonsConfis.TopRight)
	// btnTopRight := devices.NewButton(gpioTopRight)
	// gpioBottomLeft := hardware.NewGPIOSwitch(buttonsConfis.BottomLeft)
	// btnBottomLeft := devices.NewButton(gpioBottomLeft)
	// gpioBottomRight := hardware.NewGPIOSwitch(buttonsConfis.BottomRight)
	// btnBottomRight := devices.NewButton(gpioBottomRight)
	// var buzzerPin gpio.PinOut = gpioreg.ByName(configs.BuzzerGPIO)
	// buzzer := piezobuzzer.NewBuzzer(buzzerPin)

	// remoteControl := remotecontrol.NewRemoteControl(
	// 	radioDev,
	// 	roll, pitch, yaw, throttle,
	// 	btnFrontLeft, btnFrontRight,
	// 	btnToptLeft, btnTopRight,
	// 	btnBottomLeft, btnBottomRight,
	// 	configs.CommandPerSecond,
	// 	oled,
	// 	buzzer,
	// )

	// ctx, cancel := context.WithCancel(context.Background())
	// var waitGroup sync.WaitGroup

	// radioDev.StartTransmitter(ctx, &waitGroup)
	// remoteControl.Start(ctx, &waitGroup, cancel)
	// utils.WaitToAbortByESC(cancel)
	// waitGroup.Wait()
	// os.Exit(0)
}
