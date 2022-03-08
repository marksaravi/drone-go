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
	GetConnectionStateChannel() <-chan int
}
type pidControls interface {
	SetTargetStates(rotations models.Rotations)
	GetThrottles(throttle float64, rotations models.Rotations, dt time.Duration, isSafeStarted bool) models.Throttles
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
	pid                pidControls
	imu                imu
	esc                esc
	radio              radioReceiver
	logger             models.Logger
	settings           Settings
	isSafeStarted      bool
	connectionState    int
	connectionChanOpen bool
	timeout            time.Time
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
		pid:                pid,
		imu:                imu,
		esc:                esc,
		radio:              radio,
		logger:             logger,
		settings:           settings,
		timeout:            time.Now().Add(time.Second * 1000000),
		connectionState:    constants.CONNECTED,
		connectionChanOpen: true,
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
		var running bool = true
		var throttle float64 = 0
		fc.imu.ResetTime()
		for running || fc.connectionChanOpen || commandChanOpen {
			select {
			case <-ctx.Done():
				if running {
					fc.logger.Close()
					running = false
				}

			case flightCommands, ok := <-fc.radio.GetReceiverChannel():
				if ok {
					rotations := flightCommandsToRotations(flightCommands, fc.settings)
					throttle = flightCommandsToThrottle(flightCommands, fc.settings)
					fc.checkForSafeStart(throttle)
					fc.pid.SetTargetStates(rotations)
					fc.pid.Calibrate(flightCommands.ButtonTopRight, flightCommands.ButtonTopLeft)
				}
				commandChanOpen = ok

			case fc.connectionState, fc.connectionChanOpen = <-fc.radio.GetConnectionStateChannel():
				if fc.connectionChanOpen {
					fc.showConnectionState()
				}

			default:
				if running && commandChanOpen {
					rotations, imuDataAvailable := fc.imu.ReadRotations()
					if imuDataAvailable {
						throttles := fc.pid.GetThrottles(throttle, rotations.Rotations, rotations.ReadInterval, fc.isSafeStarted)
						fc.esc.SetThrottles(throttles, fc.isSafeStarted)
						fc.logger.Send(rotations)
					}
				}
			}
		}
	}()
}

func (fc *flightControl) showConnectionState() {
	switch fc.connectionState {
	case constants.CONNECTED:
		log.Println("Connected")
		fc.timeout = time.Now()
	case constants.WAITING_FOR_CONNECTION:
		log.Println("Waiting for Connection")
	case constants.DISCONNECTED:
		log.Println("Disconnected")
		fc.timeout = time.Now()
	}
}

func (fc *flightControl) checkForSafeStart(throttle float64) {
	if time.Since(fc.timeout) > SAFE_START_DURATION && throttle == 0 {
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
