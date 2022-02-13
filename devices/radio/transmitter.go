package radio

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioTransmitterLink interface {
	radioLink
	TransmitterOn()
	Transmit(models.Payload) error
	IsTransmitFailed(update bool) bool
}

type radioTransmitter struct {
	radiolink       radioTransmitterLink
	TransmitChannel chan models.FlightCommands

	connectionChannel  chan ConnectionState
	connectionState    ConnectionState
	lastConnectionTime time.Time
	connectionTimeout  time.Duration
}

func NewTransmitter(radiolink radioTransmitterLink, connectionTimeoutMs int) *radioTransmitter {
	return &radioTransmitter{
		TransmitChannel:   make(chan models.FlightCommands),
		connectionChannel: make(chan ConnectionState),
		radiolink:         radiolink,
		connectionState:   IDLE,
		connectionTimeout: time.Millisecond * time.Duration(connectionTimeoutMs),
	}
}

func (t *radioTransmitter) StartTransmitter(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Transmitter...")

	t.radiolink.TransmitterOn()
	t.radiolink.PowerOn()
	t.lastConnectionTime = time.Now()
	go func() {
		defer t.radiolink.PowerOff()
		defer wg.Done()
		defer log.Println("Transmitter is stopped.")

		var running bool = true
		for running {
			select {
			case <-ctx.Done():
				if running {
					running = false
				}

			case flightCommands := <-t.TransmitChannel:
				if t.radiolink.IsTransmitFailed(true) {
					fmt.Println("Transmit failed")
					t.radiolink.ClearStatus()
				}
				payload := utils.SerializeFlightCommand(flightCommands)
				t.radiolink.Transmit(payload)
				fmt.Println(payload)
			default:
			}
		}
	}()
}
