package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf24l01"
)

func setData(d []byte, v byte) {
	for i := 0; i < len(d); i++ {
		d[i] = v
	}
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	data := make([]byte, 8)
	setData(data, 0)
	r := nrf24l01.NewNRF24L01EnhancedBurst(
		0,
		0,
		"GPIO25",
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
	const DATA_PER_SECOND = 25
	counter := 0
	errCounter := 0
	for running {
		select {
		case _, ok := <-ctx.Done():
			if !ok {
				running = false
			}
		default:
			if time.Since(lastSent) >= time.Second/DATA_PER_SECOND {
				if r.IsTransmitFailed(true) {
					r.ClearStatus()
				}
				err := r.Transmit(data)
				if err != nil {
					errCounter++
					fmt.Println(err)
				} else {
					counter++
					fmt.Printf("%6d, %6d, %v\n", counter, errCounter, data)
				}
				setData(data, data[0]+1)
				lastSent = time.Now()
			}
		}
	}
}
