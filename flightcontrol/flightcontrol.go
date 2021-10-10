package flightcontrol

import (
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

type imu interface {
	ReadRotations() (models.ImuRotations, bool)
}

type pidControl interface {
	ApplyFlightCommands(flightCommands models.FlightCommands)
	ApplyRotations(rotations models.ImuRotations)
	Throttles() map[uint8]float32
}

type flightControl struct {
	pid        pidControl
	imu        imu
	throttles  chan<- models.Throttles
	onOff      chan<- bool
	command    <-chan models.FlightCommands
	connection <-chan bool
	logger     chan<- models.ImuRotations
}

func NewFlightControl(pid pidControl, imu imu, throttles chan<- models.Throttles, onOff chan<- bool, command <-chan models.FlightCommands, connection <-chan bool, logger chan<- models.ImuRotations) *flightControl {
	return &flightControl{
		pid:        pid,
		imu:        imu,
		throttles:  throttles,
		onOff:      onOff,
		command:    command,
		connection: connection,
		logger:     logger,
	}
}

func (fc *flightControl) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(fc.onOff)
		defer close(fc.throttles)
		defer close(fc.logger)
		defer log.Println("Flight Control stopped")
		fc.onOff <- true
		for fc.command != nil || fc.connection != nil {
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
			default:
			}
		}
	}()
}
