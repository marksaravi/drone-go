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

type KalmanFilter interface {
	Value(value float64) float64
}

type DroneSettings struct {
	ImuMems            imuMems
	Receiver           radioReceiver
	Escs               escs
	ImuDataPerSecond   int
	ESCsDataPerSecond  int
	CommandsPerSecond  int
	PlotterActive      bool
	RollMidValue       int
	PitchMidValue      int
	YawMidValue        int
	RotationRange      float64
	MaxThrottle        float64
	ThrottleZeroOffset int
	PID                pid.PIDSettings
}

type droneApp struct {
	startTime          time.Time
	imuDataPerSecond   int
	imu                imuMems
	escs               escs
	escsDataPerImuData int
	flightControl      *FlightControl

	commandsPerSecond int
	receiver          radioReceiver
	lastImuRead       time.Time
	imuReadInterval   time.Duration
	plotterActive     bool

	rollMidValue       int
	pitchlMidValue     int
	yawMidValue        int
	rotationRange      float64
	maxThrottle        float64
	throttleZeroOffset int
	throttle           KalmanFilter
}

func NewDrone(settings DroneSettings) *droneApp {
	escsDataPerImuData := settings.ImuDataPerSecond / settings.ESCsDataPerSecond
	return &droneApp{
		startTime:          time.Now(),
		imu:                settings.ImuMems,
		escs:               settings.Escs,
		escsDataPerImuData: escsDataPerImuData,
		flightControl:      NewFlightControl(settings.Escs, settings.MaxThrottle, settings.PID, escsDataPerImuData),
		imuDataPerSecond:   settings.ImuDataPerSecond,
		receiver:           settings.Receiver,
		commandsPerSecond:  settings.CommandsPerSecond,
		lastImuRead:        time.Now(),
		imuReadInterval:    time.Second / time.Duration(settings.ImuDataPerSecond),
		plotterActive:      settings.PlotterActive,
		rollMidValue:       settings.RollMidValue,
		pitchlMidValue:     settings.PitchMidValue,
		yawMidValue:        settings.YawMidValue,
		rotationRange:      settings.RotationRange,
		maxThrottle:        settings.MaxThrottle,
		throttleZeroOffset: settings.ThrottleZeroOffset,
		throttle:           utils.NewKalmanFilter(0.15, 1),
	}
}

func (d *droneApp) Start(ctx context.Context, wg *sync.WaitGroup) {
	var commandOk, running bool = true, true
	var commands []byte

	fmt.Println("Starting Drone...")
	d.flightControl.turnOnMotors(false)

	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)
	imuReadTick := utils.WithDataPerSecond(d.imuDataPerSecond)
	escCounter := utils.NewCounter(d.escsDataPerImuData)

	for running || commandOk {
		select {
		case commands, commandOk = <-commandsChannel:
			if commandOk {
				d.applyCommands(commands)
			}
		case _, running = <-ctx.Done():
			running = false
		default:
			if imuReadTick.IsTime() {
				rot, err := d.imu.Read()
				if err == nil {
					d.flightControl.calcOutputThrottles(rot)
				}
			}
			if escCounter.Inc() {
				go func() {
					d.flightControl.applyThrottles()
				}()
			}
		}
	}
	fmt.Println("Stopping Motors...")
	d.flightControl.turnOnMotors(false)
	fmt.Println("Stopping Drone...")
}
