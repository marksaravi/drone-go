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

	rollMidValue   int
	pitchlMidValue int
	yawMidValue    int
	rotationRange  float64
	maxThrottle    float64
	throttle       KalmanFilter
}

func NewDrone(settings DroneSettings) *droneApp {
	escsDataPerImuData := settings.ImuDataPerSecond / settings.ESCsDataPerSecond
	return &droneApp{
		startTime:          time.Now(),
		imu:                settings.ImuMems,
		escs:               settings.Escs,
		escsDataPerImuData: escsDataPerImuData,
		flightControl: NewFlightControl(
			settings.Escs,
			settings.MaxThrottle,
			settings.MaxOutputThrottle,
			settings.Arm_0_2_Pid,
			settings.Arm_1_3_Pid,
			settings.Yaw_Pid,
			escsDataPerImuData,
		),
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
		throttle:          utils.NewKalmanFilter(0.15, 1),
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
					if escCounter.Inc() {
						d.flightControl.applyThrottles()
					}
				}
			}
		}
	}
	fmt.Println("Stopping Motors...")
	d.flightControl.turnOnMotors(false)
	fmt.Println("Stopping Drone...")
}
