package main

import (
	"context"
	"fmt"
	"log"

	dronePackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	log.Println("Starting RemoteControl")
	configs := dronePackage.ReadConfigs("./configs/drone-configs.json")
	log.Println(configs)

	radioConfigs := configs.Radio
	radioLink := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI.BusNum,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.SPI.SpiChipEnabledGPIO,
		radioConfigs.RxTxAddress,
	)
	radioReceiver := radio.NewReceiver(radioLink)

	ctx, cancel := context.WithCancel(context.Background())
	drone := dronePackage.NewDrone(dronePackage.DroneSettings{
		Receiver:          radioReceiver,
		CommandsPerSecond: configs.CommandsPerSecond,
	})

	go func() {
		fmt.Scanln()
		cancel()
	}()

	drone.Start(ctx)
}
