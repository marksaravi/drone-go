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

type radioReceiverLink interface {
	radioLink
	ReceiverOn()
	Listen()
	Receive() (models.Payload, error)
	IsReceiverDataReady(update bool) bool
}

type radioReceiver struct {
	radiolink           radioReceiverLink
	receiveChannel      chan models.FlightCommands
	commandReadInterval time.Duration
	commandReadTimeout  time.Time

	connectionChannel  chan int
	connectionState    int
	lastConnectionTime time.Time
	connectionTimeout  time.Duration
}

func NewReceiver(radiolink radioReceiverLink, commandsPerSecond int, connectionTimeoutMs int) *radioReceiver {
	return &radioReceiver{
		receiveChannel:      make(chan models.FlightCommands),
		connectionChannel:   make(chan int),
		radiolink:           radiolink,
		connectionState:     constants.IDLE,
		commandReadInterval: time.Second / time.Duration(commandsPerSecond*2),
		commandReadTimeout:  time.Now(),
		connectionTimeout:   time.Millisecond * time.Duration(connectionTimeoutMs),
		lastConnectionTime:  time.Now(),
	}
}

func (r *radioReceiver) StartReceiver(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Receiver...")

	r.radiolink.ReceiverOn()
	r.radiolink.PowerOn()
	r.radiolink.Listen()

	flushTimeout := time.Now()
	go func() {
		defer r.radiolink.PowerOff()
		defer wg.Done()
		defer log.Println("Receiver is stopped.")

		for {
			select {
			case <-ctx.Done():
				close(r.connectionChannel)
				close(r.receiveChannel)
				return

			default:
				if time.Since(r.commandReadTimeout) >= r.commandReadInterval {
					r.commandReadTimeout = time.Now()
					if r.radiolink.IsReceiverDataReady(true) {
						payload, _ := r.radiolink.Receive()
						r.radiolink.Listen()
						if time.Since(flushTimeout) > time.Second {
							r.receiveChannel <- utils.DeserializeFlightCommand(payload)
							r.updateConnectionState(true)
						}
					} else {
						r.updateConnectionState(false)
					}
				}
			}
		}
	}()
}
func (r *radioReceiver) updateConnectionStateAsync(connectionState int) {
	r.connectionState = connectionState
	func() {
		r.connectionChannel <- connectionState
	}()
}
func (r *radioReceiver) updateConnectionState(connected bool) {
	if connected {
		r.lastConnectionTime = time.Now()
		if r.connectionState != constants.CONNECTED {
			r.updateConnectionStateAsync(constants.CONNECTED)
		}
	}
	if !connected && r.connectionState != constants.DISCONNECTED && time.Since(r.lastConnectionTime) > r.connectionTimeout {
		r.updateConnectionStateAsync(constants.DISCONNECTED)
	}
}

func (r *radioReceiver) GetReceiverChannel() <-chan models.FlightCommands {
	return r.receiveChannel
}

func (r *radioReceiver) GetConnectionStateChannel() <-chan int {
	return r.connectionChannel
}
