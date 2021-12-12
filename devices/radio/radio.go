package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type radioDevice struct {
	transmitter      chan models.FlightCommands
	receiver         chan models.FlightCommands
	connection       chan bool
	radio            models.RadioLink
	connected        bool
	lastHeartBeat    time.Time
	heartBeatTimeout time.Duration
}

func NewRadio(radio models.RadioLink, heartBeatTimeoutMs int) *radioDevice {
	return &radioDevice{
		transmitter:      make(chan models.FlightCommands),
		receiver:         make(chan models.FlightCommands),
		connection:       make(chan bool),
		radio:            radio,
		heartBeatTimeout: time.Duration(heartBeatTimeoutMs * int(time.Millisecond)),
		lastHeartBeat:    time.Now(),
		connected:        false,
	}
}

func (r *radioDevice) Acknowledge() bool {
	return r.Transmit(models.FlightCommands{
		Id:   0,
		Time: time.Now().UnixNano(),
	})
}

func (r *radioDevice) Transmit(data models.FlightCommands) bool {
	if r.transmitter != nil {
		r.transmitter <- data
		return true
	}
	return false
}

func (r *radioDevice) Close() {
	close(r.transmitter)
	r.transmitter = nil
}

func (r *radioDevice) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Radio...")

	go func() {
		defer wg.Done()
		defer log.Println("Radio is stopped...")

		r.setConnection(false)
		r.radio.ReceiverOn()
		var running bool = true
		for running && r.transmitter != nil {
			select {
			case <-ctx.Done():
				close(r.connection)
				close(r.receiver)
				running = false
			case data, ok := <-r.transmitter:
				if ok {
					r.radio.TransmitterOn()
					r.radio.Transmit(utils.SerializeFlightCommand(data))
					r.radio.ReceiverOn()
				}
			default:
				data, available := r.radio.Receive()
				if available {
					r.receiver <- utils.DeserializeFlightCommand(data)
				}
				r.setConnection(available)
			}
		}
	}()
}

func (r *radioDevice) setConnection(available bool) {
	if available {
		if !r.connected {
			r.connected = true
			r.connection <- true
		}
		r.lastHeartBeat = time.Now()
	} else {
		if r.connected {
			if time.Since(r.lastHeartBeat) > r.heartBeatTimeout {
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
