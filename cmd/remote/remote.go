package main

import (
	"context"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/apps/remote"
	pushbutton "github.com/marksaravi/drone-go/devices/push-button"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
	"github.com/marksaravi/drone-go/hardware/ssd1306"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting RemoteControl")
	configs := remote.ReadConfigs("./configs/remote-configs.json")
	log.Println(configs)

	radioConfigs := configs.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI,
		radioConfigs.RxTxAddress,
	)

	radioTransmitter := radio.NewRadioTransmitter(radioLink)

	analogToDigitalSPIConn := hardware.NewSPIConnection(configs.Joysticks.SPI)

	joystickRoll := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.RollChannel,
	)
	joystickPitch := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.PitchChannel,
	)
	joystickYaw := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.YawChannel,
	)
	joystickThrottle := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		configs.Joysticks.ThrottleChannel,
	)

	buttons := make([]remote.PushButton, 0, 10)
	buttonsCount := make([]int, 0, 10)
	for i := 0; i < len(configs.PushButtons); i++ {
		pin := hardware.NewPushButtonInput(configs.PushButtons[i].GPIO)
		buttons = append(buttons, pushbutton.NewPushButton(configs.PushButtons[i].Name, pin, configs.PushButtons[i].PulseMode))
		buttonsCount = append(buttonsCount, 0)
	}

	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()
	d := &i2c.Dev{Addr: 0x3D, Bus: b}
	oled := ssd1306.NewSSD1306(d, ssd1306.DefaultOptions)
	err = oled.Init()
	if err != nil {
		log.Fatal(err)
	}
	remoteControl := remote.NewRemoteControl(remote.RemoteSettings{
		Transmitter:            radioTransmitter,
		CommandPerSecond:       configs.CommandsPerSecond,
		Roll:                   joystickRoll,
		Pitch:                  joystickPitch,
		Yaw:                    joystickYaw,
		Throttle:               joystickThrottle,
		PushButtons:            buttons,
		OLED:                   oled,
		DisplayUpdatePerSecond: configs.DisplayUpdatePerSecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Scanln()
		cancel()
	}()

	remoteControl.Start(ctx)
}
