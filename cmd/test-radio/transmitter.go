package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/models"
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
	)
	transmitter := radio.NewTransmitter(radioNRF204)
	go transmitter.StartTransmitter(ctx, wg)

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		ts := time.Now()
		var throttle uint16 = 0
		for {
			select {
			case <-ctx.Done():
				close(transmitter.TransmitChannel)
				return
			default:
				if time.Since(ts) >= time.Second/time.Duration(configs.CommandPerSecond) {
					ts = time.Now()
					transmitter.TransmitChannel <- models.FlightCommands{Throttle: throttle}
					throttle++
				}
			}
		}
	}(ctx, wg)
}
