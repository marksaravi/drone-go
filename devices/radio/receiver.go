package radio

import (
	"context"
	"log"
	"sync"
	"time"

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

	r.radiolink.ReceiverOn()

	go func() {
		defer r.radiolink.PowerOff()
		defer wg.Done()
		defer log.Println("Receiver is stopped.")

		ts := time.Now()
		for {
			select {
			case <-ctx.Done():
				return

			default:
				if time.Since(ts) >= time.Second/10 {
					ts = time.Now()
					if r.radiolink.ReceiverDataReady(true) {

					}
				}
			}
		}
	}()
}
