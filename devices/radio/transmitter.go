package radio

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

func NewTransmitter(radiolink radioLink) *radioDevice {
	return &radioDevice{
		transmitter:     make(chan models.FlightCommands),
		receiver:        make(chan models.FlightCommands),
		connection:      make(chan ConnectionState),
		radiolink:       radiolink,
		connectionState: IDLE,
	}
}

func (r *radioDevice) StartTransmitter(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Transmitter...")
	r.radiolink.SetTransmitterAddress()

	go func() {
		defer wg.Done()
		defer log.Println("Transmitter is stopped.")

		var running bool = true
		for running {
			select {
			case <-ctx.Done():
				if running {
					running = false
				}

			default:
			}
		}
	}()
}
