package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/constants"
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
	SetThrottles(models.Throttles)
}

type pidControls interface {
	SetPIDTargetState(state models.PIDState)
	SetRotations(rotations models.ImuRotations)
	Throttles() map[int]float64
	PrintGains()
}

type Settings struct {
	MaxThrottle float64
	MaxRoll     float64
	MaxPitch    float64
	MaxYaw      float64
}

type flightControl struct {
	pid      pidControls
	imu      imu
	esc      esc
	radio    models.Radio
	logger   models.Logger
	settings Settings
}

func NewFlightControl(
	pid pidControls,
	imu imu,
	esc esc,
	radio models.Radio,
	logger models.Logger,
	settings Settings,
) *flightControl {
	return &flightControl{
		pid:      pid,
		imu:      imu,
		esc:      esc,
		radio:    radio,
		logger:   logger,
		settings: settings,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("Flight Control is stopped...")
		defer fc.pid.PrintGains()
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
					fc.radio.Close()
					fc.logger.Close()
					running = false
				}

			case flightCommands, ok := <-fc.radio.GetReceiver():
				if ok {
					fc.pid.SetPIDTargetState(flightCommandsToPIDState(flightCommands, fc.settings))
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
						fc.esc.SetThrottles(models.Throttles{
							Active:    true,
							Throttles: fc.pid.Throttles(),
						})
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

func joystickToTwoWayCommand(digital uint16, resolution uint16, max float64) float64 {
	return (float64(digital) - float64(resolution/2)) / float64(resolution) * max
}

func joystickToOneWayCommand(digital uint16, resolution uint16, max float64) float64 {
	return float64(digital) / float64(resolution) * max
}

func flightCommandsToPIDState(command models.FlightCommands, settings Settings) models.PIDState {
	return models.PIDState{
		Roll:     joystickToTwoWayCommand(command.Roll, constants.JOYSTICK_RESOLUTION, settings.MaxRoll),
		Pitch:    joystickToTwoWayCommand(command.Pitch, constants.JOYSTICK_RESOLUTION, settings.MaxPitch),
		Yaw:      joystickToTwoWayCommand(command.Yaw, constants.JOYSTICK_RESOLUTION, settings.MaxYaw),
		Throttle: joystickToOneWayCommand(command.Throttle, constants.JOYSTICK_RESOLUTION, settings.MaxThrottle),
	}
}
