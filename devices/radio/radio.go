package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type ConnectionState = int

const (
	CONNECTED ConnectionState = iota
	DISCONNECTED
	LOST
)

const (
	DATA_PAYLOAD byte = iota
	HEART_BEAT_PAYLOAD
	RECEIVER_OFF
)

type radioDevice struct {
	transmitter           chan models.FlightCommands
	receiver              chan models.FlightCommands
	connection            chan ConnectionState
	radio                 models.RadioLink
	connectionState       ConnectionState
	lastSentHeartBeat     time.Time
	lastReceivedHeartBeat time.Time
	heartBeatTimeout      time.Duration
}

func NewRadio(radio models.RadioLink, heartBeatTimeoutMs int) *radioDevice {
	heartBeatTimeout := time.Duration(heartBeatTimeoutMs * int(time.Millisecond))
	hearBeatInit := time.Now().Add(-heartBeatTimeout * 2)
	return &radioDevice{
		transmitter:           make(chan models.FlightCommands),
		receiver:              make(chan models.FlightCommands),
		connection:            make(chan ConnectionState),
		radio:                 radio,
		heartBeatTimeout:      heartBeatTimeout,
		connectionState:       DISCONNECTED,
		lastSentHeartBeat:     hearBeatInit,
		lastReceivedHeartBeat: hearBeatInit,
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
		for running || transmitterChannelOpen {

			if running {
				payload, available := r.radio.Receive()
				flightCommands := utils.DeserializeFlightCommand(payload)
				if available && flightCommands.PayloadType != HEART_BEAT_PAYLOAD {
					r.receiver <- flightCommands
				}
				r.setConnection(available)
			}

			select {
			case flightCommands, ok := <-r.transmitter:
				transmitterChannelOpen = ok
				if transmitterChannelOpen {
					r.transmit(flightCommands)
				}
			default:
				if transmitterChannelOpen && time.Since(r.lastSentHeartBeat) >= r.heartBeatTimeout/2 {
					r.transmit(models.FlightCommands{
						Id:          0,
						PayloadType: HEART_BEAT_PAYLOAD,
						Time:        time.Now().UnixNano(),
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
	r.transmitter <- data
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
		if r.connected {
			r.connection <- CONNECTED
		} else {
			r.connection <- DISCONNECTED
		}

	}
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan int {
	return r.connection
}
