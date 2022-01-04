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
	IDLE ConnectionState = iota
	DISCONNECTED
	CONNECTED
	CONNECTION_LOST
)

const (
	NO_COMMAND models.FlightCommandType = iota
	COMMAND
	HEARTBEAT
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
		connectionState:       IDLE,
		lastSentHeartBeat:     hearBeatInit,
		lastReceivedHeartBeat: hearBeatInit,
	}
}

func (r *radioDevice) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Radio...")

	go func() {
		defer wg.Done()
		defer log.Println("Radio is stopped...")

		r.clearBuffer()

		for {
			select {
			case <-ctx.Done():
				r.closeRadio()
				log.Println("Closing receiver and connection...")
				close(r.receiver)
				close(r.connection)
				return

			case flightCommands, ok := <-r.transmitter:
				if ok {
					r.lastSentHeartBeat = time.Now()
					r.radio.Transmit(utils.SerializeFlightCommand(flightCommands))
				}

			default:
				payload, available := r.radio.Receive()
				if payload[0] == COMMAND {
					r.receiver <- utils.DeserializeFlightCommand(payload)
				}
				if available {
					r.setConnectionState(payload[0])
				} else {
					r.setConnectionState(NO_COMMAND)
				}
			}

			if time.Since(r.lastSentHeartBeat) >= r.heartBeatTimeout/4 {
				r.radio.Transmit(utils.SerializeFlightCommand(models.FlightCommands{
					Type: HEARTBEAT,
				}))
				r.lastSentHeartBeat = time.Now()
			}
		}
	}()
}

func (r *radioDevice) closeRadio() {
	var receiverOffPayload models.Payload
	receiverOffPayload[0] = RECEIVER_OFF
	for i := 0; i < 5; i++ {
		r.radio.Transmit(receiverOffPayload)
		time.Sleep(time.Millisecond)
	}
}

func (r *radioDevice) setConnectionState(commandType models.FlightCommandType) {
	prevState := r.connectionState
	switch commandType {
	case HEARTBEAT, COMMAND:
		r.connectionState = CONNECTED
		r.lastReceivedHeartBeat = time.Now()
	case NO_COMMAND:
		if time.Since(r.lastReceivedHeartBeat) > r.heartBeatTimeout && r.connectionState == CONNECTED {
			r.connectionState = CONNECTION_LOST
		}
		if r.connectionState == IDLE {
			r.connectionState = DISCONNECTED
		}
	case RECEIVER_OFF:
		r.connectionState = DISCONNECTED
	}
	if prevState != r.connectionState {
		r.connection <- r.connectionState
	}
}

func (r *radioDevice) clearBuffer() {
	for {
		_, available := r.radio.Receive()
		if !available {
			break
		}
	}
	log.Println("Radio buffer is cleared.")
}

func (r *radioDevice) Transmit(data models.FlightCommands) {
	data.Type = COMMAND
	r.transmitter <- data
}

func (r *radioDevice) Close() {
	close(r.transmitter)
}

func (r *radioDevice) SuppressLostConnection() {
	// if r.connectionState == DISCONNECTED {
	// 	return
	// }
	// r.actions <- SUPPRESS_LOST_CONNECTION
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan ConnectionState {
	return r.connection
}
