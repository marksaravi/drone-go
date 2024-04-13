package drone

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/pid"
	"github.com/marksaravi/drone-go/utils"
)

type radioReceiver interface {
	Start(ctx context.Context, wg *sync.WaitGroup, commandsPerSecond int) <-chan []byte
}

type imuMems interface {
	Read() (imu.Rotations, error)
}

type escs interface {
	On()
	Off()
	SetThrottles(motors []float64)
}

type DroneSettings struct {
	ImuMems           imuMems
	Receiver          radioReceiver
	Escs              escs
	ImuDataPerSecond  int
	CommandsPerSecond int
	PlotterActive     bool
	RollMidValue      int
	PitchMidValue     int
	YawMidValue       int
	RotationRange     float64
	MaxThrottle       float64
	MinFlightThrottle float64
	PID               pid.PIDSettings
}

type droneApp struct {
	startTime        time.Time
	imuDataPerSecond int
	imu              imuMems
	escs             escs
	flightControl    *FlightControl

	commandsPerSecond int
	receiver          radioReceiver
	lastImuRead       time.Time
	imuReadInterval   time.Duration
	plotterActive     bool

	rollMidValue      int
	pitchlMidValue    int
	yawMidValue       int
	rotationRange     float64
	maxThrottle       float64
	minFlightThrottle float64
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		startTime:         time.Now(),
		imu:               settings.ImuMems,
		escs:              settings.Escs,
		flightControl:     NewFlightControl(settings.Escs, settings.MinFlightThrottle, settings.MaxThrottle, settings.PID),
		imuDataPerSecond:  settings.ImuDataPerSecond,
		receiver:          settings.Receiver,
		commandsPerSecond: settings.CommandsPerSecond,
		lastImuRead:       time.Now(),
		imuReadInterval:   time.Second / time.Duration(settings.ImuDataPerSecond),
		plotterActive:     settings.PlotterActive,
		rollMidValue:      settings.RollMidValue,
		pitchlMidValue:    settings.PitchMidValue,
		yawMidValue:       settings.YawMidValue,
		rotationRange:     settings.RotationRange,
		maxThrottle:       settings.MaxThrottle,
		minFlightThrottle: settings.MinFlightThrottle,
	}
}

func (d *droneApp) Start(ctx context.Context, wg *sync.WaitGroup) {
	var commandOk, running bool = true, true
	var commands []byte

	fmt.Println("Starting Drone...")
	d.flightControl.turnOnMotors(false)

	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)
	imuReadTick := utils.WithDataPerSecond(d.imuDataPerSecond)
	escUpdates := utils.WithDataPerSecond(40)
	for running || commandOk {
		select {
		case commands, commandOk = <-commandsChannel:
			if commandOk {
				go func() {
					d.applyCommands(commands)
				}()
			}
		case _, running = <-ctx.Done():
			running = false
		default:
			if imuReadTick.IsTime() {
				rot, err := d.imu.Read()
				if err == nil {
					d.flightControl.SetRotations(rot)
				}
			}
			if escUpdates.IsTime() {
				go func() {
					d.flightControl.SetMotorsPowers()
				}()
			}
		}
	}
	fmt.Println("Stopping Motors...")
	d.flightControl.turnOnMotors(false)
	fmt.Println("Stopping Drone...")
}
