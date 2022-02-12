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
	TransmitPayload(models.Payload) error
	// ReceivePayload() (models.Payload, bool)
}

func process(ctx context.Context, wg *sync.WaitGroup, radio radioLink) {
	defer wg.Done()
	wg.Add(1)

	interval := time.Second / 40
	start := time.Now()
	var received int = 0
	lastReceived := time.Now()

	var running bool = true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
		}
		if time.Since(start) >= interval {
			start = time.Now()
			var payload models.Payload
			radio.TransmitPayload(payload)
		}
		_, ok := radio.ReceivePayload()
		if ok {
			received++
		}
		if time.Since(lastReceived) >= time.Second {
			log.Println(received)
			received = 0
			lastReceived = time.Now()
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
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go process(ctx, &wg, nrf204dev)
	utils.WaitToAbortByENTER(cancel, &wg)
	wg.Wait()
}
