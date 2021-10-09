package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

type imu interface {
	ReadRotations() (models.ImuRotations, bool)
}

type flightControl struct {
	imu        imu
	command    <-chan models.FlightCommands
	connection <-chan bool
	logger     chan<- models.ImuRotations
}

func NewFlightControl(imu imu, command <-chan models.FlightCommands, connection <-chan bool, logger chan<- models.ImuRotations) *flightControl {
	return &flightControl{
		imu:        imu,
		command:    command,
		connection: connection,
		logger:     logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("Flight Control stopped")
		for fc.command != nil || fc.connection != nil || fc.logger != nil {
			rotations, imuDataAvailable := fc.imu.ReadRotations()
			if imuDataAvailable {
				if fc.logger != nil {
					fc.logger <- rotations
				}
			}
			select {
			case _, isCommandOk := <-fc.command:
				if !isCommandOk {
					fc.command = nil
				}
			case cnonnected, isConnectionOk := <-fc.connection:
				if isConnectionOk {
					log.Println("Connected: ", cnonnected)
				} else {
					fc.connection = nil
				}
			case <-ctx.Done():
				if fc.logger != nil {
					close(fc.logger)
					fc.logger = nil
				}
			default:
			}
		}
	}()
}
