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
	data := make([]byte, 8)

	r := nrf24l01.NewNRF24L01EnhancedBurst(
		0,
		0,
		"GPIO24",
		"03896",
	)
	r.TransmitterOn()
	r.PowerOn()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Scanln()
		cancel()
	}()
	lastSent := time.Now()
	running := true
	for running {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				running = false
			}
		default:
			if time.Since(lastSent) >= time.Second/2 {
				err := r.Transmit(data)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(data)
				data[0] += 1
				lastSent = time.Now()
			}
		}
	}
}
