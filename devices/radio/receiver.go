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
	statusCheckInterval time.Duration

	connectionChannel  chan models.ConnectionState
	connectionState    models.ConnectionState
	lastConnectionTime time.Time
	connectionTimeout  time.Duration
}

func NewReceiver(radiolink radioReceiverLink, commandsPerSecond int, connectionTimeoutMs int) *radioReceiver {
	return &radioReceiver{
		receiveChannel:      make(chan models.FlightCommands),
		connectionChannel:   make(chan models.ConnectionState),
		radiolink:           radiolink,
		connectionState:     constants.WAITING_FOR_CONNECTION,
		statusCheckInterval: time.Second / time.Duration(commandsPerSecond*2),
		connectionTimeout:   time.Millisecond * time.Duration(connectionTimeoutMs),
	}
}

func (r *radioReceiver) StartReceiver(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Receiver...")

	r.radiolink.ReceiverOn()
	r.radiolink.PowerOn()
	r.radiolink.Listen()
	r.lastConnectionTime = time.Now()
	go func() {
		defer r.radiolink.PowerOff()
		defer wg.Done()
		defer log.Println("Receiver is stopped.")

		ts := time.Now()
		for {
			select {
			case <-ctx.Done():
				close(r.connectionChannel)
				close(r.receiveChannel)
				return

			default:
				if time.Since(ts) >= r.statusCheckInterval {
					ts = time.Now()
					if r.radiolink.IsReceiverDataReady(true) {
						payload, _ := r.radiolink.Receive()
						r.radiolink.Listen()
						r.receiveChannel <- utils.DeserializeFlightCommand(payload)
					}
				}
			}
		}
	}()
}

func (r *radioReceiver) GetReceiverChannel() <-chan models.FlightCommands {
	return r.receiveChannel
}

func (r *radioReceiver) GetConnectionStateChannel() <-chan models.ConnectionState {
	return r.connectionChannel
}
