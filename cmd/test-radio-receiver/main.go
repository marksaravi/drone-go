package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()

	r := nrf24l01.NewNRF24L01EnhancedBurst(
		0,
		0,
		"GPIO24",
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
	lastSent := time.Now()
	running := true
	update := false
	for running {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				running = false
			}
		default:
			if time.Since(lastSent) >= time.Second/5 {
				lastSent = time.Now()
				if r.IsReceiverDataReady(update) {
					data, err := r.Receive()
					r.Listen()
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(data)
					}
					fmt.Println("Listening...")
				}
				update = !update
			}
		}
	}

}
