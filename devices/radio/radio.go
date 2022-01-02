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
	CONNECTED
	DISCONNECTED
	LOST
)

const (
	NO_PAYLOAD byte = iota
	DATA_PAYLOAD
	HEARTBEAT_PAYLOAD
	RECEIVER_OFF_PAYLOAD
)

type radioDevice struct {
	transmitter           chan models.FlightCommands
	receiver              chan models.FlightCommands
	connection            chan ConnectionState
	radio                 models.RadioLink
	connectionState       ConnectionState
	ConnectionStateLock   sync.Mutex
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

func (r *radioDevice) transmitPayload(payload models.Payload) {
	r.lastSentHeartBeat = time.Now()
	r.radio.TransmitterOn()
	r.radio.Transmit(payload)
	r.radio.ReceiverOn()
}

func (r *radioDevice) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Radio...")

	go func() {
		defer wg.Done()
		defer log.Println("Radio is stopped...")

		r.clearBuffer()
		var running bool = true
		var transmitterOpen bool = true
		for running || transmitterOpen {
			select {
			case <-ctx.Done():
				if running {
					log.Println("Closing receiver and connection")
					for i := 0; i < 3; i++ {
						r.transmitPayload(genPayload(RECEIVER_OFF_PAYLOAD))
						time.Sleep(time.Millisecond * 50)
					}
					time.Sleep(time.Millisecond * 100)
					close(r.receiver)
					r.receiver = nil
					close(r.connection)
					r.connection = nil
					running = false
				}

			case flightCommands, ok := <-r.transmitter:
				if ok {
					r.transmitPayload(utils.SerializeFlightCommand(flightCommands))
				}
				transmitterOpen = ok

			default:
			}

			if running {
				r.receivePayload()
				r.sendHeartbeat()
			}
		}
	}()
}

func (r *radioDevice) receivePayload() {
	payload, available := r.radio.Receive()
	if available {
		payloadType := payload[0]
		r.setConnectionState(available, payloadType)
		if payloadType == DATA_PAYLOAD && r.receiver != nil {
			flightCommands := utils.DeserializeFlightCommand(payload)
			r.receiver <- flightCommands
		}
	} else {
		r.setConnectionState(false, NO_PAYLOAD)
	}
}

func (r *radioDevice) sendHeartbeat() {
	if time.Since(r.lastSentHeartBeat) >= r.heartBeatTimeout/4 {
		r.transmitPayload(genPayload(HEARTBEAT_PAYLOAD))
	}
}

func (r *radioDevice) setConnectionState(available bool, payloadType byte) {
	r.ConnectionStateLock.Lock()
	defer r.ConnectionStateLock.Unlock()

	prevState := r.connectionState
	if available {
		r.connectionState = CONNECTED
		r.lastReceivedHeartBeat = time.Now()
		if payloadType == RECEIVER_OFF_PAYLOAD {
			r.connectionState = DISCONNECTED
		}
	} else {
		if time.Since(r.lastReceivedHeartBeat) > r.heartBeatTimeout && r.connectionState == CONNECTED {
			r.connectionState = LOST
		}
	}
	if prevState != r.connectionState && r.connection != nil {
		r.connection <- r.connectionState
	}
}

func (r *radioDevice) SuppressLostConnection() {
	go func() {
		log.Println("suppressing lost connection...")
		r.setConnectionState(true, RECEIVER_OFF_PAYLOAD)
	}()
}
func (r *radioDevice) clearBuffer() {
	for {
		_, available := r.radio.Receive()
		if !available {
			break
		}
	}
}

func (r *radioDevice) Transmit(data models.FlightCommands) {
	data.PayloadType = DATA_PAYLOAD
	if r.transmitter != nil {
		r.transmitter <- data
	}
}

func (r *radioDevice) CloseTransmitter() {
	close(r.transmitter)
	r.transmitter = nil
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan int {
	return r.connection
}

func genPayload(payloadType byte) models.Payload {
	var payload models.Payload
	payload[0] = payloadType
	return models.Payload(payload)
}
