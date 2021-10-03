package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

type flightControl struct {
	imu        <-chan models.ImuRotations
	command    <-chan models.FlightCommands
	connection <-chan bool
	logger     chan<- models.ImuRotations
}

func NewFlightControl(imu <-chan models.ImuRotations, command <-chan models.FlightCommands, connection <-chan bool, logger chan<- models.ImuRotations) *flightControl {
	return &flightControl{
		imu:        imu,
		command:    command,
		connection: connection,
		logger:     logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer log.Println("Flight Control stopped")
	for {
		select {
		case rotations := <-fc.imu:
			fc.logger <- rotations
		case <-fc.command:
		case cnonnected := <-fc.connection:
			log.Println("Connected: ", cnonnected)
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}
