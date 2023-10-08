package main

import (
	"context"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()

	r := nrf24l01.NewNRF24L01EnhancedBurst(
		hardware.SPIConnConfigs{
			BusNumber:       1,
			ChipSelect:      0,
			ChipEnabledGPIO: "GPIO16",
		},
		"03896",
	)
	r.ReceiverOn()
	r.PowerOn()
	r.Listen()
	fmt.Println("Listening...")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Scanln()
		cancel()
	}()

	counter := 0
	errCounter := 0
	running := true
	update := false
	for running {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				running = false
			}
		default:
			if r.IsReceiverDataReady(true) {
				data, err := r.Receive()
				r.Listen()
				if err != nil {
					errCounter++
					fmt.Println(err)
				} else {
					counter++
					fmt.Printf("%6d, %6d, %v\n", counter, errCounter, data)
				}
				fmt.Println("Listening...")
			}
			update = !update
		}
	}

}
