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
	transmitChannel chan models.FlightCommands

	connectionChannel  chan int
	connectionState    int
	lastConnectionTime time.Time
	connectionTimeout  time.Duration
}

func NewTransmitter(radiolink radioTransmitterLink, connectionTimeoutMs int) *radioTransmitter {
	return &radioTransmitter{
		transmitChannel:   make(chan models.FlightCommands),
		connectionChannel: make(chan int),
		radiolink:         radiolink,
		connectionState:   constants.IDLE,
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
				close(t.connectionChannel)
				return

			case flightCommands := <-t.transmitChannel:
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

func (t *radioTransmitter) Close() {
	close(t.transmitChannel)
}

func (t *radioTransmitter) GetConnectionStateChannel() <-chan int {
	return t.connectionChannel
}

func (t *radioTransmitter) SuppressLostConnection() {
	t.connectionState = constants.IDLE
}

func (t *radioTransmitter) Transmit(flightCommands models.FlightCommands) {
	t.transmitChannel <- flightCommands
}

func (t *radioTransmitter) updateConnectionState(connected bool) {
	prevState := t.connectionState
	t.connectionState, t.lastConnectionTime = newConnectionState(connected, t.connectionState, t.lastConnectionTime, t.connectionTimeout)
	if prevState != t.connectionState {
		t.connectionChannel <- t.connectionState
	}
}

func newConnectionState(
	connected bool,
	prevState int,
	lastConnected time.Time,
	timeout time.Duration,
) (newState int, lastConnection time.Time) {
	newState = prevState
	lastConnection = lastConnected
	if connected {
		newState = constants.CONNECTED
		lastConnection = time.Now()
	} else {
		if prevState == constants.IDLE {
			newState = constants.WAITING_FOR_CONNECTION
		} else if prevState == constants.CONNECTED {
			if time.Since(lastConnected) < timeout {
				newState = constants.CONNECTED
			} else {
				newState = constants.DISCONNECTED
			}
		}
	}
	return
}
