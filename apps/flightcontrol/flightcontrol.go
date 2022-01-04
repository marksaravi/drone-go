package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/models"
)

type imu interface {
	ResetTime()
	ReadRotations() (models.ImuRotations, bool)
}

type esc interface {
	On()
	Off()
	Close()
	SetThrottles(models.Throttles)
}

type pidControls interface {
	SetFlightCommands(flightCommands models.FlightCommands)
	SetRotations(rotations models.ImuRotations)
	Throttles() models.Throttles
	SetEmergencyStop(stop bool)
}

type flightControl struct {
	pid    pidControls
	imu    imu
	esc    esc
	radio  models.Radio
	logger models.Logger
}

func NewFlightControl(
	pid pidControls,
	imu imu,
	esc esc,
	radio models.Radio,
	logger models.Logger,
) *flightControl {
	return &flightControl{
		pid:    pid,
		imu:    imu,
		esc:    esc,
		radio:  radio,
		logger: logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("Flight Control is stopped...")
		defer fc.esc.Off()

		fc.esc.On()
		var commandChanOpen bool = true
		var connectionChanOpen bool = true
		var running bool = true

		fc.imu.ResetTime()
		for running || connectionChanOpen || commandChanOpen {
			select {
			case <-ctx.Done():
				if running {
					fc.esc.Close()
					fc.radio.Close()
					fc.logger.Close()
					running = false
				}

			case flightCommands, ok := <-fc.radio.GetReceiver():
				if ok {
					fc.pid.SetFlightCommands(flightCommands)
					// showFlightCommands(flightCommands)
				}
				commandChanOpen = ok

			case connectionState, ok := <-fc.radio.GetConnection():
				if ok {
					showConnectionState(connectionState)
				}
				connectionChanOpen = ok

			default:
				if running && commandChanOpen {
					rotations, imuDataAvailable := fc.imu.ReadRotations()
					if imuDataAvailable {
						fc.pid.SetRotations(rotations)
						fc.esc.SetThrottles(fc.pid.Throttles())
						fc.logger.Send(rotations)
					}
				}
			}
		}
	}()
}

func showConnectionState(connectionState radio.ConnectionState) {
	switch connectionState {
	case radio.CONNECTED:
		log.Println("Connected")
	case radio.DISCONNECTED:
		log.Println("Disconnected")
	case radio.CONNECTION_LOST:
		log.Println("Lost")
	}
}

// var lastShowFlightCommands time.Time

// func showFlightCommands(fc models.FlightCommands) {
// 	if time.Since(lastShowFlightCommands) >= time.Second/2 {
// 		lastShowFlightCommands = time.Now()
// 		log.Printf("%4d, %4d, %4d, %4d, %t, %t, %t, %t, %t, %t", fc.Roll, fc.Pitch, fc.Yaw, fc.Throttle, fc.ButtonFrontLeft, fc.ButtonFrontRight, fc.ButtonTopLeft, fc.ButtonTopRight, fc.ButtonBottomLeft, fc.ButtonBottomRight)
// 	}
// }
