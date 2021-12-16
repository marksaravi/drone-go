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
		// var lastPrinted time.Time = time.Now()
		// var running bool = true
		// command := fc.radio.GetReceiver()
		// connection := fc.radio.GetConnection()
		// for running {
		// 	select {
		// 	case <-ctx.Done():
		// 		fc.radio.Close()
		// 		fc.logger.Close()
		// 		running = false
		// 	case flightCommands, ok := <-command:
		// 		if ok {
		// 			if time.Since(lastPrinted) >= time.Second {
		// 				showFLightCommands(flightCommands)
		// 				lastPrinted = time.Now()
		// 			}
		// 		} else {
		// 			command = nil
		// 		}
		// 	case connected, ok := <-connection:
		// 		if ok {
		// 			log.Println("connected: ", connected)
		// 		} else {
		// 			log.Println("channel is closed")
		// 			connection = nil
		// 		}
		// 	default:
		// 		rotations, imuDataAvailable := fc.imu.ReadRotations()
		// 		if imuDataAvailable {
		// 			if fc.logger != nil {
		// 				fc.logger.Send(rotations)
		// 			}
		// 		}
		// 	}
		// }
	}()
}

func showFLightCommands(fc models.FlightCommands) {
	log.Printf("%8.2f, %8.2f, %t, %t", fc.Roll, fc.Pitch, fc.ButtonFrontLeft, fc.ButtonTopLeft)
}
