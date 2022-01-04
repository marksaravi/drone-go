package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func process(ctx context.Context, wg *sync.WaitGroup, radio models.RadioLink) {
	defer wg.Done()

	wg.Add(1)
	transmitInterval := time.Second / 50
	lastTransmit := time.Now()
	var id byte = 0
	var running bool = true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			if time.Since(lastTransmit) >= transmitInterval {
				lastTransmit = time.Now()
				var payload models.Payload
				payload[0] = id
				id++
				if id > 250 {
					id = 0
				}
				if id%50 == 0 {
					log.Println("50 sent")
				}
				err := radio.Transmit(payload)
				if err != nil {
					log.Println(err.Error())
				}
			}

		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	radioConfigs := config.ReadConfigs().RemoteControl.Radio
	radioNRF204 := nrf204.NewRadio(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go process(ctx, &wg, radioNRF204)
	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
