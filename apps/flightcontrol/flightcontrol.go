package flightcontrol

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/models"
)

const SAFE_START_DURATION = time.Second

type imu interface {
	ResetTime()
	ReadRotations() (models.ImuRotations, bool)
}

type esc interface {
	On()
	Off()
	SetThrottles(throttles models.Throttles, isSafeStarted bool)
}

type radioReceiver interface {
	GetReceiverChannel() <-chan models.FlightCommands
	GetConnectionStateChannel() <-chan models.ConnectionState
}
type pidControls interface {
	SetTargetStates(rotations models.Rotations, throttle float64)
	SetStates(rotations models.Rotations, dt time.Duration)
	GetThrottles(isSafeStarted bool) models.Throttles

	PrintGains()
	Calibrate(up, down bool)
}

type Settings struct {
	MaxThrottle float64
	MaxRoll     float64
	MaxPitch    float64
	MaxYaw      float64
}

type flightControl struct {
	pid           pidControls
	imu           imu
	esc           esc
	radio         radioReceiver
	logger        models.Logger
	settings      Settings
	isSafeStarted bool
}

func NewFlightControl(
	pid pidControls,
	imu imu,
	esc esc,
	radio radioReceiver,
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
		var flightControlStartTime time.Time = time.Now()
		fc.imu.ResetTime()
		for running || connectionChanOpen || commandChanOpen {
			select {
			case <-ctx.Done():
				if running {
					fc.logger.Close()
					running = false
				}

			case flightCommands, ok := <-fc.radio.GetReceiverChannel():
				if ok {
					rotations := flightCommandsToRotations(flightCommands, fc.settings)
					throttle := flightCommandsToThrottle(flightCommands, fc.settings)
					fc.checkForSafeStart(throttle, flightControlStartTime)
					fc.pid.SetTargetStates(rotations, throttle)
					fc.pid.Calibrate(flightCommands.ButtonTopRight, flightCommands.ButtonTopLeft)
				}
				commandChanOpen = ok

			case connectionState, ok := <-fc.radio.GetConnectionStateChannel():
				if ok {
					showConnectionState(connectionState)
				}
				connectionChanOpen = ok

			default:
				if running && commandChanOpen {
					rotations, imuDataAvailable := fc.imu.ReadRotations()
					if imuDataAvailable {
						fc.pid.SetStates(rotations.Rotations, rotations.ReadInterval)
						throttles := fc.pid.GetThrottles(fc.isSafeStarted)
						fc.esc.SetThrottles(throttles, fc.isSafeStarted)
						fc.logger.Send(rotations)
					}
				}
			}
		}
	}()
}

func showConnectionState(connectionState models.ConnectionState) {
	switch connectionState {
	case constants.CONNECTED:
		log.Println("Connected")
	case constants.WAITING_FOR_CONNECTION:
		log.Println("Waiting for Connection")
	case constants.DISCONNECTED:
		log.Println("Disconnected")
	}
}

func (fc *flightControl) checkForSafeStart(throttle float64, startTime time.Time) {
	if time.Since(startTime) > SAFE_START_DURATION && throttle == 0 {
		if !fc.isSafeStarted {
			log.Println("Safe Start Detected")
		}
		fc.isSafeStarted = true
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

func flightCommandsToRotations(command models.FlightCommands, settings Settings) models.Rotations {
	return models.Rotations{
		Roll:  joystickToTwoWayCommand(command.Roll, constants.JOYSTICK_RESOLUTION, settings.MaxRoll),
		Pitch: joystickToTwoWayCommand(command.Pitch, constants.JOYSTICK_RESOLUTION, settings.MaxPitch),
		Yaw:   joystickToTwoWayCommand(command.Yaw, constants.JOYSTICK_RESOLUTION, settings.MaxYaw),
	}
}

func flightCommandsToThrottle(command models.FlightCommands, settings Settings) float64 {
	return joystickToOneWayCommand(command.Throttle, constants.JOYSTICK_RESOLUTION, settings.MaxThrottle)
}
