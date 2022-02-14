package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware/nrf204"
)

func runReceiver(ctx context.Context, wg *sync.WaitGroup) {
	configs := config.ReadConfigs().FlightControl
	radioConfigs := configs.Radio
	log.Println(radioConfigs)

	radioNRF204 := nrf204.NewNRF204EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)
	receiver := radio.NewReceiver(radioNRF204, configs.CommandPerSecond, radioConfigs.ConnectionTimeoutMs)
	go receiver.StartReceiver(ctx, wg)
}
