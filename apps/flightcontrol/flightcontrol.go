package flightcontrol

import (
	"context"
	"fmt"
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
	SetThrottles(throttles models.Throttles)
}

type radioReceiver interface {
	GetReceiverChannel() <-chan models.FlightCommands
	GetConnectionStateChannel() <-chan int
}
type pidControls interface {
	SetTargetStates(rotations models.Rotations)
	GetThrottles(throttle float64, rotations models.Rotations, dt time.Duration) models.Throttles
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
	commandChanOpen    bool
	connectionChanOpen bool
	running            bool
	flightCommands     models.FlightCommands
	throttle           float64
	stopTimeout        time.Time
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
		connectionState:    constants.CONNECTED,
		connectionChanOpen: true,
		commandChanOpen:    true,
		running:            true,
		throttle:           0,
		isSafeStarted:      false,
		stopTimeout:        time.Now(),
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
		fc.imu.ResetTime()
		for fc.running || fc.connectionChanOpen || fc.commandChanOpen {
			select {
			case <-ctx.Done():
				if fc.running {
					fc.logger.Close()
					fc.running = false
				}

			case fc.flightCommands, fc.commandChanOpen = <-fc.radio.GetReceiverChannel():
				if fc.commandChanOpen {
					rotations := flightCommandsToRotations(fc.flightCommands, fc.settings)
					throttle := flightCommandsToThrottle(fc.flightCommands, fc.settings)
					fc.checkForEnablingSafeStart(throttle)
					fc.pid.SetTargetStates(rotations)
					fc.pid.Calibrate(fc.flightCommands.ButtonTopRight, fc.flightCommands.ButtonTopLeft)
				}

			case fc.connectionState, fc.connectionChanOpen = <-fc.radio.GetConnectionStateChannel():
				if fc.connectionChanOpen {
					fc.isSafeStarted = false
					fc.showConnectionState()
				}

			default:
				fc.safeReduceThrottle()
				if fc.running && fc.commandChanOpen {
					rotations, imuDataAvailable := fc.imu.ReadRotations()
					if imuDataAvailable {
						throttles := fc.pid.GetThrottles(fc.throttle, rotations.Rotations, rotations.ReadInterval)
						fc.esc.SetThrottles(throttles)
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
	case constants.DISCONNECTED:
		log.Println("Disconnected")
	}
}

func (fc *flightControl) safeReduceThrottle() {
	if fc.isSafeStarted || fc.throttle == 0 || time.Since(fc.stopTimeout) < time.Millisecond*10 {
		return
	}
	fc.stopTimeout = time.Now()
	fc.throttle -= fc.throttle / 100
	if fc.throttle < 5 {
		fc.throttle = 0
	}
}

func (fc *flightControl) checkForEnablingSafeStart(throttle float64) {
	if !fc.isSafeStarted && throttle == 0 {
		fc.isSafeStarted = true
		fmt.Println("Safe Start Enabled")
	}
	if fc.isSafeStarted {
		fc.throttle = throttle
	}
}

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
