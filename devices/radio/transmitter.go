package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/constants"
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
	t.radiolink.Transmit(utils.SerializeFlightCommand(models.FlightCommands{
		Type: constants.COMMAND_DUMMY,
	}))
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
				t.radiolink.Transmit(utils.SerializeFlightCommand(flightCommands))
			}
		}
	}()
}

func (t *radioTransmitter) updateConnectionState(connected bool) {
	prevState := t.connectionState
	t.connectionState, t.lastConnectionTime = newConnectionState(connected, t.connectionState, t.lastConnectionTime, t.connectionTimeout)
	if prevState != t.connectionState {
		t.ConnectionChannel <- t.connectionState
	}
}

func newConnectionState(
	connected bool,
	prevState ConnectionState,
	lastConnected time.Time,
	timeout time.Duration,
) (newState ConnectionState, lastConnection time.Time) {
	newState = prevState
	lastConnection = lastConnected
	if connected {
		newState = CONNECTED
		lastConnection = time.Now()
	} else {
		if prevState == IDLE {
			newState = WAITING_FOR_CONNECTION
		} else if prevState == CONNECTED {
			if time.Since(lastConnected) < timeout {
				newState = CONNECTED
			} else {
				newState = DISCONNECTED
			}
		}
	}
	return
}
