package radio

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

func NewReceiver(radiolink radioLink) *radioReceiver {
	return &radioReceiver{
		receiveChannel:    make(chan models.FlightCommands),
		connectionChannel: make(chan ConnectionState),
		radiolink:         radiolink,
		connectionState:   IDLE,
	}
}

func (r *radioReceiver) StartReceiver(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Starting the Receiver...")

	go func() {
		defer wg.Done()
		defer log.Println("Receiver is stopped.")

		var running bool = true
		for running {
			select {
			case <-ctx.Done():
				if running {
					running = false
				}

			default:
			}
		}
	}()
}
