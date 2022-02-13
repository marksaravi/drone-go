package radio

import (
	"context"
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

	ConnectionChannel  chan ConnectionState
	connectionState    ConnectionState
	lastConnectionTime time.Time
	connectionTimeout  time.Duration
}

func NewTransmitter(radiolink radioTransmitterLink, connectionTimeoutMs int) *radioTransmitter {
	return &radioTransmitter{
		TransmitChannel:   make(chan models.FlightCommands),
		ConnectionChannel: make(chan ConnectionState),
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

		for {
			select {
			case <-ctx.Done():
				close(t.ConnectionChannel)
				return

			case flightCommands := <-t.TransmitChannel:
				if t.radiolink.IsTransmitFailed(true) {
					t.updateConnectionState(false)
					t.radiolink.ClearStatus()
				} else {
					t.updateConnectionState(true)
				}
				payload := utils.SerializeFlightCommand(flightCommands)
				t.radiolink.Transmit(payload)
			}
		}
	}()
}

func (t *radioTransmitter) updateConnectionState(connected bool) {
	prevState := t.connectionState
	if connected {
		t.connectionState = CONNECTED
		t.lastConnectionTime = time.Now()
	} else {
		if t.connectionState == IDLE {
			t.lastConnectionTime = time.Now()
			return
		}
		if t.connectionState == CONNECTED && time.Since(t.lastConnectionTime) > t.connectionTimeout {
			t.connectionState = DISCONNECTED
		}
	}
	if prevState != t.connectionState {
		t.ConnectionChannel <- t.connectionState
	}
}
