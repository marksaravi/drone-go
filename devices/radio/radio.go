package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
)

type radioDevice struct {
	transmitter chan models.FlightCommands
	receiver    chan models.FlightCommands
	connection  chan bool
	radio       models.RadioLink
}

func NewRadio(radio models.RadioLink) *radioDevice {
	return &radioDevice{
		transmitter: make(chan models.FlightCommands),
		receiver:    make(chan models.FlightCommands),
		connection:  make(chan bool),
		radio:       radio,
	}
}

func (r *radioDevice) Acknowledge() {
	r.Transmit(models.FlightCommands{
		Id:   0,
		Time: time.Now().UnixNano(),
	})
}

func (r *radioDevice) Transmit(data models.FlightCommands) bool {
	if r.transmitter != nil {
		r.transmitter <- data
		return true
	}
	return false
}

func (r *radioDevice) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Radio...")

	go func() {
		defer wg.Done()
		defer log.Println("Radio is stopped...")

		// var connected bool = false
		// r.connection <- connected
		var running bool = true
		for running {
			select {
			case <-ctx.Done():
				close(r.transmitter)
				r.transmitter = nil
				close(r.connection)
				close(r.receiver)
				log.Println("Stopping the Radio...")
				running = false
			default:
				// data, available := rl.Receive()
				// if available {
				// 	log.Println(utils.DeserializeFlightCommand(data))
				// 	if !connected {
				// 		connected = true
				// 		r.connection <- connected
				// 	}
				// }
			}
		}
	}()
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan bool {
	return r.connection
}
