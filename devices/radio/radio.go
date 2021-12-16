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
	return &radioDevice{
		transmitter:           make(chan models.FlightCommands),
		receiver:              make(chan models.FlightCommands),
		connection:            make(chan bool),
		radio:                 radio,
		heartBeatTimeout:      time.Duration(heartBeatTimeoutMs * int(time.Millisecond)),
		connected:             false,
		lastSentHeartBeat:     time.Now(),
		lastReceivedHeartBeat: time.Now(),
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
		var done bool = false
		for !done {
			select {
			case <-ctx.Done():
				if !done {
					log.Println("RADIO CONTEXT DONE")
					done = true
				}

			default:
			}
		}

		// r.setConnection(false)
		// r.radio.ReceiverOn()
		// var running bool = true
		// for running && r.transmitter != nil {
		// 	select {
		// 	case <-ctx.Done():
		// 		close(r.connection)
		// 		close(r.receiver)
		// 		running = false
		// 	case data, ok := <-r.transmitter:
		// 		if ok {
		// 			r.transmit(data)
		// 		}
		// 	default:
		// 		payload, available := r.radio.Receive()
		// 		if available {
		// 			data := utils.DeserializeFlightCommand(payload)
		// 			if data.Type == DATA_PAYLOAD {
		// 				r.receiver <- data
		// 			}
		// 		}
		// 		r.setConnection(available)
		// 		if time.Since(r.lastSentHeartBeat) >= r.heartBeatTimeout/2 {
		// 			r.transmit(models.FlightCommands{
		// 				Id:   0,
		// 				Type: HEART_BEAT_PAYLOAD,
		// 				Time: time.Now().UnixNano(),
		// 			})
		// 		}
		// 	}
		// }
	}()
}

func (r *radioDevice) Transmit(data models.FlightCommands) bool {
	if r.transmitter != nil {
		r.transmitter <- data
		return true
	}
	return false
}

func (r *radioDevice) Close() {
	log.Println("Closing...")
	close(r.transmitter)
	r.transmitter = nil
}

func (r *radioDevice) setConnection(available bool) {
	if available {
		if !r.connected {
			r.connected = true
			r.connection <- true
		}
		r.lastReceivedHeartBeat = time.Now()
	} else {
		if r.connected {
			if time.Since(r.lastReceivedHeartBeat) > r.heartBeatTimeout {
				r.connected = false
				r.connection <- false
			}
		}
	}
}

func (r *radioDevice) GetReceiver() <-chan models.FlightCommands {
	return r.receiver
}

func (r *radioDevice) GetConnection() <-chan bool {
	return r.connection
}
