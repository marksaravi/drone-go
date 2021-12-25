package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type imu interface {
	ResetTime()
	ReadRotations() (models.ImuRotations, bool)
}

type pidControl interface {
	ApplyFlightCommands(flightCommands models.FlightCommands)
	ApplyRotations(rotations models.ImuRotations)
	Throttles() map[uint8]float32
}

type flightControl struct {
	pid    pidControl
	imu    imu
	radio  models.Radio
	logger models.Logger
}

func NewFlightControl(
	pid pidControl,
	imu imu,
	radio models.Radio,
	logger models.Logger,
) *flightControl {
	return &flightControl{
		pid:    pid,
		imu:    imu,
		radio:  radio,
		logger: logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("Flight Control is stopped...")

		var flightCommands models.FlightCommands
		var connectionChanOpen bool = true
		var connectionState radio.ConnectionState = radio.DISCONNECTED
		var receiverChanOpen bool = true
		var running bool = true
		fc.imu.ResetTime()
		for running || connectionChanOpen || receiverChanOpen {
			rotations, imuDataAvailable := fc.imu.ReadRotations()
			if running && imuDataAvailable {
				fc.logger.Send(rotations)
			}

			select {
			case connectionState, connectionChanOpen = <-fc.radio.GetConnection():
				if connectionChanOpen {
					log.Println("Connected: ", connectionState)
				}
			default:
			}

			select {
			case flightCommands, receiverChanOpen = <-fc.radio.GetReceiver():
				if receiverChanOpen {
					utils.SerializeFlightCommand(flightCommands)
				}
			default:
			}

			select {
			case <-ctx.Done():
				if running {
					fc.radio.CloseTransmitter()
					fc.logger.Close()
					running = false
				}
			default:
			}
		}
	}()
}

func showFLightCommands(fc models.FlightCommands) {
	log.Printf("%8.2f, %8.2f, %t, %t", fc.Roll, fc.Pitch, fc.ButtonFrontLeft, fc.ButtonTopLeft)
}
