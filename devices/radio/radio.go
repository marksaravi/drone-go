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

type radioLink interface {
	TransmitPayload(models.Payload) error
	ReceivePayload() (models.Payload, bool)
}

type radioDevice struct {
	transmitter           chan models.FlightCommands
	receiver              chan models.FlightCommands
	connection            chan ConnectionState
	radiolink             radioLink
	connectionState       ConnectionState
	lastSentHeartBeat     time.Time
	lastReceivedHeartBeat time.Time
	heartBeatTimeout      time.Duration
	lock                  sync.Mutex
}

func NewRadio(radiolink radioLink, heartBeatTimeoutMs int) *radioDevice {
	heartBeatTimeout := time.Duration(heartBeatTimeoutMs * int(time.Millisecond))
	hearBeatInit := time.Now().Add(-heartBeatTimeout * 2)
	return &radioDevice{
		transmitter:           make(chan models.FlightCommands),
		receiver:              make(chan models.FlightCommands),
		connection:            make(chan ConnectionState),
		radiolink:             radiolink,
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
		defer log.Println("Radio is stopped.")

		var heartbeatInterval = r.heartBeatTimeout / 10
		r.clearBuffer()
		var running bool = true
		var transmitting bool = true
		for running || transmitting {
			select {
			case <-ctx.Done():
				if running {
					r.closeRadio()
					log.Println("Closing Receiver and Connection...")
					close(r.receiver)
					close(r.connection)
					running = false
				}

			case flightCommands, ok := <-r.transmitter:
				if ok {
					r.transmitPayload(utils.SerializeFlightCommand(flightCommands))
				}
				transmitting = ok
			default:
				if running {
					payload, available := r.radiolink.ReceivePayload()
					if payload[0] == COMMAND {
						r.receiver <- utils.DeserializeFlightCommand(payload)
					}
					if available {
						r.setConnectionState(payload[0])
					} else {
						r.setConnectionState(NO_COMMAND)
					}
					if time.Since(r.lastSentHeartBeat) >= heartbeatInterval {
						r.transmitPayload(utils.SerializeFlightCommand(models.FlightCommands{
							Type: HEARTBEAT,
						}))
					}
				}
			}
		}
	}()
}

func (r *radioDevice) Transmit(data models.FlightCommands) {
	go func() {
		data.Type = COMMAND
		r.transmitter <- data
	}()
}

func (r *radioDevice) Close() {
	log.Println("Closing Transmitter...")
	go func() {
		close(r.transmitter)
	}()
}

func (r *radioDevice) SuppressLostConnection() {
	go func() {
		r.setConnectionState(RECEIVER_OFF)
	}()
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan ConnectionState {
	return r.connection
}

func (r *radioDevice) closeRadio() {
	var receiverOffPayload models.Payload
	receiverOffPayload[0] = RECEIVER_OFF
	for i := 0; i < 5; i++ {
		r.transmitPayload(receiverOffPayload)
	}
}

func (r *radioDevice) transmitPayload(payload models.Payload) {
	r.lastSentHeartBeat = time.Now()
	r.radiolink.TransmitPayload(payload)
}

func (r *radioDevice) setConnectionState(commandType models.FlightCommandType) {
	defer r.lock.Unlock()

	r.lock.Lock()
	prevState := r.connectionState
	switch commandType {
	case NO_COMMAND:
		if time.Since(r.lastReceivedHeartBeat) > r.heartBeatTimeout && r.connectionState == CONNECTED {
			r.connectionState = CONNECTION_LOST
		}
		if r.connectionState == IDLE {
			r.connectionState = DISCONNECTED
		}
	case HEARTBEAT, COMMAND:
		r.connectionState = CONNECTED
		r.lastReceivedHeartBeat = time.Now()
	case RECEIVER_OFF:
		r.connectionState = DISCONNECTED
	}
	if prevState != r.connectionState {
		r.connection <- r.connectionState
	}
}

func (r *radioDevice) clearBuffer() {
	for {
		_, available := r.radiolink.ReceivePayload()
		if !available {
			break
		}
	}
	log.Println("Radio buffer is cleared.")
}
