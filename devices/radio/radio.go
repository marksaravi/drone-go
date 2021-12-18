package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

const (
	DATA_PAYLOAD byte = iota
	HEART_BEAT_PAYLOAD
)

type radioDevice struct {
	transmitter           chan models.FlightCommands
	receiver              chan models.FlightCommands
	connection            chan bool
	radio                 models.RadioLink
	connected             bool
	lastSentHeartBeat     time.Time
	lastReceivedHeartBeat time.Time
	heartBeatTimeout      time.Duration
}

func NewRadio(radio models.RadioLink, heartBeatTimeoutMs int) *radioDevice {
	heartBeatTimeout := time.Duration(heartBeatTimeoutMs * int(time.Millisecond))
	return &radioDevice{
		transmitter:           make(chan models.FlightCommands),
		receiver:              make(chan models.FlightCommands),
		connection:            make(chan bool),
		radio:                 radio,
		heartBeatTimeout:      heartBeatTimeout,
		connected:             false,
		lastSentHeartBeat:     time.Now(),
		lastReceivedHeartBeat: time.Now().Add(-heartBeatTimeout * 2),
	}
}

func (r *radioDevice) transmit(data models.FlightCommands) {
	r.lastSentHeartBeat = time.Now()
	r.radio.TransmitterOn()
	r.radio.Transmit(utils.SerializeFlightCommand(data))
	r.radio.ReceiverOn()
}

func (r *radioDevice) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Radio...")

	go func() {
		defer wg.Done()
		defer log.Println("Radio is stopped...")
		var running bool = true
		var transmitterChannelOpen bool = true
		var flightCommands models.FlightCommands
		for running || transmitterChannelOpen {

			if running {
				payload, available := r.radio.Receive()
				if available {
					r.receiver <- utils.DeserializeFlightCommand(payload)
				}
				r.setConnection(available)
			}

			select {
			case flightCommands, transmitterChannelOpen = <-r.transmitter:
				if transmitterChannelOpen {
					r.transmit(flightCommands)
				}
			default:
				if transmitterChannelOpen && time.Since(r.lastSentHeartBeat) >= r.heartBeatTimeout/2 {
					r.transmit(models.FlightCommands{
						Id:   0,
						Type: HEART_BEAT_PAYLOAD,
						Time: time.Now().UnixNano(),
					})
				}
			}

			select {
			case <-ctx.Done():
				if running {
					log.Println("Closing receiver and connection")
					close(r.receiver)
					close(r.connection)
					running = false
				}
			default:
			}

		}
	}()
}

func (r *radioDevice) Transmit(data models.FlightCommands) {
	// r.transmitter <- data
}

func (r *radioDevice) CloseTransmitter() {
	close(r.transmitter)
}

func (r *radioDevice) setConnection(available bool) {
	connected := false
	if available {
		r.lastReceivedHeartBeat = time.Now()
	}
	if time.Since(r.lastReceivedHeartBeat) < r.heartBeatTimeout {
		connected = true
	}
	if connected != r.connected {
		r.connected = connected
		r.connection <- connected
	}
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan bool {
	return r.connection
}
