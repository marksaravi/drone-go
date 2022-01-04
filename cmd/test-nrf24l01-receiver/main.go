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

type radioLink interface {
	Receive() (models.Payload, bool)
	Transmit(models.Payload) error
}

func process(ctx context.Context, wg *sync.WaitGroup, radio radioLink) {
	defer wg.Done()
	wg.Add(1)

	start := time.Now()

	var total int = 0
	var counter int = 0
	var failed int = 0
	var id byte = 0

	var running bool = true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			payload, available := radio.Receive()
			if available {
				total++
				counter++
				if id != payload[0] {
					failed++
				}
				id = payload[0] + 1
				if id > 250 {
					id = 0
				}
			}
			if time.Since(start) >= time.Second {
				log.Println("   Data Per Second: ", counter, total, failed)
				counter = 0
				start = time.Now()
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
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go process(ctx, &wg, nrf204dev)
	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
