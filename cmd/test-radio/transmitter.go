package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware/nrf204"
)

func runTransmitter(ctx context.Context, wg *sync.WaitGroup) {
	configs := config.ReadConfigs().RemoteControl
	radioConfigs := configs.Radio
	log.Println(radioConfigs)

	radioNRF204 := nrf204.NewNRF204EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	transmitter := radio.NewTransmitter(radioNRF204)
	go transmitter.StartTransmitter(ctx, wg)
}
