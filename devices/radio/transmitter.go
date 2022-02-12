package radio

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func NewTransmitter(radiolink radioLink) *radioTransmitter {
	return &radioTransmitter{
		transmitChannel:   make(chan models.FlightCommands),
		connectionChannel: make(chan ConnectionState),
		radiolink:         radiolink,
		connectionState:   IDLE,
	}
}

func (r *radioTransmitter) StartTransmitter(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Transmitter...")
	r.radiolink.TransmitterOn()

	go func() {
		defer r.radiolink.PowerOff()
		defer wg.Done()
		defer log.Println("Transmitter is stopped.")

		var running bool = true
		for running {
			select {
			case <-ctx.Done():
				if running {
					running = false
				}

			case flightCommands := <-r.transmitChannel:
				payload := utils.SerializeFlightCommand(flightCommands)
				r.radiolink.Transmit(payload)
			default:
			}
		}
	}()
}
