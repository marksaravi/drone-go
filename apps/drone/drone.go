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
	ReadAll() (imu.Rotations, imu.Rotations, imu.Rotations, error)
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
	ImuMems           imuMems
	Receiver          radioReceiver
	Escs              escs
	ImuDataPerSecond  int
	ESCsDataPerSecond int
	CommandsPerSecond int
	PlotterActive     bool
	RollMidValue      int
	PitchMidValue     int
	YawMidValue       int
	RotationRange     float64
	MaxThrottle       float64
	MaxOutputThrottle float64
	Arm_0_2_Pid       *pid.PIDControl
	Arm_1_3_Pid       *pid.PIDControl
	Yaw_Pid           *pid.PIDControl
	RollDirection     float64
	PitchDirection    float64
}

type droneApp struct {
	startTime          time.Time
	imuDataPerSecond   int
	escUpdatePerSecond int
	imu                imuMems
	escs               escs
	flightControl      *FlightControl
	commandsPerSecond  int
	receiver           radioReceiver
	lastImuRead        time.Time
	imuReadInterval    time.Duration
	plotterActive      bool

	rollMidValue   int
	pitchlMidValue int
	yawMidValue    int
	rotationRange  float64
	maxThrottle    float64
	throttle       KalmanFilter
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		startTime: time.Now(),
		imu:       settings.ImuMems,
		escs:      settings.Escs,
		flightControl: NewFlightControl(
			settings.Escs,
			settings.MaxThrottle,
			settings.MaxOutputThrottle,
			settings.Arm_0_2_Pid,
			settings.Arm_1_3_Pid,
			settings.Yaw_Pid,
			settings.RollDirection,
			settings.PitchDirection,
		),
		imuDataPerSecond:   settings.ImuDataPerSecond,
		escUpdatePerSecond: settings.ESCsDataPerSecond,
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
		throttle:           utils.NewKalmanFilter(0.15, 1),
	}
}

func (d *droneApp) Start(ctx context.Context, wg *sync.WaitGroup) {
	var commandOk, running bool = true, true
	var commands []byte

	fmt.Println("Starting Drone...")
	d.flightControl.turnOnMotors(false)
	// ******** IMU DUR:  1.966206139s 19.662µs
	// ******** ESCs DUR:  52.049043ms 520ns
	// ******** COMMAND DUR:  2.511909837s 25.119µs
	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)
	d.lastImuRead = time.Now()
	imuReadCycles := d.imuDataPerSecond / d.escUpdatePerSecond
	escUpdateCounter := 0

	for running || commandOk {
		select {
		case commands, commandOk = <-commandsChannel:
			if commandOk {
				d.applyCommands(commands)
			}
		case _, running = <-ctx.Done():
			running = false
		default:
			if time.Since(d.lastImuRead) >= d.imuReadInterval {
				d.lastImuRead = time.Now()
				rot, _, grot, _ := d.imu.ReadAll()
				d.flightControl.calcOutputThrottles(rot, grot)
				escUpdateCounter++
				if escUpdateCounter >= imuReadCycles {
					d.flightControl.applyThrottles()
					escUpdateCounter = 0
				}
			}
		}
	}
	fmt.Println("Stopping Motors...")
	d.flightControl.turnOnMotors(false)
	fmt.Println("Stopping Drone...")
}
