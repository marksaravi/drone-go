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
	log.Println(configs)
	radioConfigs := configs.Radio

	radioNRF204 := nrf204.NewNRF204(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	transmitter := radio.NewTransmitter(radioNRF204, radioConfigs.HeartBeatTimeoutMS)
	go transmitter.StartTransmitter(ctx, wg)
}
