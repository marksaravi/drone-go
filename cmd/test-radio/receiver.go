package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func runReceiver(ctx context.Context, wg *sync.WaitGroup) {
	configs := config.ReadConfigs().FlightControl
	radioConfigs := configs.Radio
	log.Println(radioConfigs)

	radioNRF24L01 := nrf24l01.NewNRF24L01EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)
	receiver := radio.NewReceiver(radioNRF24L01, configs.CommandPerSecond, radioConfigs.ConnectionTimeoutMs)
	go receiver.Start(ctx, wg)
}
