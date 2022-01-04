package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func process(ctx context.Context, wg *sync.WaitGroup, radiodev models.Radio) {
	defer wg.Done()
	wg.Add(1)

	interval := time.Second / 400
	start := time.Now()
	var id byte = 0

	var running bool = true
	for running {
		select {
		case <-ctx.Done():
			return
		case connection, ok := <-radiodev.GetConnection():
			if ok {
				var label string
				switch connection {
				case radio.IDLE:
					label = "IDLE"
				case radio.CONNECTED:
					label = "CONNECTED"
				case radio.DISCONNECTED:
					label = "DISCONNECTED"
				case radio.CONNECTION_LOST:
					label = "CONNECTION_LOST"
				}
				log.Println("CONNECTION: ", label)
			}
		case data, ok := <-radiodev.GetReceiver():
			if ok {
				log.Println("DATA: ", data)
			}
		default:
			if time.Since(start) >= interval {
				start = time.Now()
				var payload models.Payload
				payload[0] = id
				radiodev.Transmit(models.FlightCommands{
					Type: id,
				})
				id++
				if id > 250 {
					id = 0
				}
			}
		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	radioConfigs := config.ReadConfigs().RemoteControl.Radio
	nrf204dev := nrf204.NewNRF204(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	radio := radio.NewRadio(nrf204dev, radioConfigs.HeartBeatTimeoutMS)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	radio.Start(&wg)
	go process(ctx, &wg, radio)
	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
