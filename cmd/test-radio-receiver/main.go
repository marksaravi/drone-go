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

type radioLink interface {
	Receive() (models.Payload, bool)
	Transmit(models.Payload) error
}

func process(ctx context.Context, wg *sync.WaitGroup, radiodev models.Radio) {
	defer wg.Done()
	wg.Add(1)

	var counter int = 0
	var start time.Time = time.Now()

	for {
		select {
		case <-ctx.Done():
			radiodev.Close()
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
		case _, ok := <-radiodev.GetReceiver():
			if ok {
				counter++
				if time.Since(start) >= time.Second {
					log.Println("DATA PER SECOND: ", counter)
					start = time.Now()
					counter = 0
				}
			}
		}
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	radioConfigs := config.ReadConfigs().FlightControl.Radio
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
	radio.Start(ctx, &wg)
	go process(ctx, &wg, radio)
	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
