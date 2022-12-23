package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/marksaravi/drone-go/apps/remotecontrol"
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
	piezobuzzer "github.com/marksaravi/drone-go/hardware/piezo-buzzer"
	"github.com/marksaravi/drone-go/hardware/ssd1306"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	log.Println("Starting RemoteControl")
	configs := config.ReadConfigs().RemoteControl
	log.Println(configs)
	radioConfigs := configs.Radio
	joysticksConfigs := configs.Joysticks
	buttonsConfis := configs.Buttons
	hardware.HostInitialize()

	oledConn, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer oledConn.Close()
	oledDev := &i2c.Dev{Addr: configs.DisplayAddress, Bus: oledConn}
	oled := ssd1306.NewSSD1306(oledDev, ssd1306.DefaultOptions)
	oled.Init()

	radioNRF24L01 := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)
	radioDev := radio.NewTransmitter(radioNRF24L01, radioConfigs.ConnectionTimeoutMs)
	analogToDigitalSPIConn := hardware.NewSPIConnection(
		joysticksConfigs.SPI.BusNumber,
		joysticksConfigs.SPI.ChipSelect,
	)
	xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.Roll.Channel,
	)
	yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.Pitch.Channel,
	)
	zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.Yaw.Channel,
	)
	throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		joysticksConfigs.Throttle.Channel,
	)
	roll := devices.NewJoystick(xAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Roll.Offset, joysticksConfigs.Roll.Dir)
	pitch := devices.NewJoystick(yAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Pitch.Offset, joysticksConfigs.Pitch.Dir)
	yaw := devices.NewJoystick(zAxisAnalogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Yaw.Offset, joysticksConfigs.Yaw.Dir)
	throttle := devices.NewJoystick(throttleAlogToDigitalConvertor, int(mcp3008.DIGITAL_MAX_VALUE), joysticksConfigs.Throttle.Offset, 1)
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
		oled,
		buzzer,
	)

	ctx, cancel := context.WithCancel(context.Background())
	var waitGroup sync.WaitGroup

	radioDev.StartTransmitter(ctx, &waitGroup)
	remoteControl.Start(ctx, &waitGroup, cancel)
	utils.WaitToAbortByESC(cancel)
	waitGroup.Wait()
	os.Exit(0)
}
