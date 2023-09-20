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
}
