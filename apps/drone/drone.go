package drone

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/apps/plotter"
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

const PLOTTER_ADDRESS = "192.168.1.101:8000"

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
	PID               PIDConfigs
}

type droneApp struct {
	startTime        time.Time
	imuDataPerSecond int
	imu              imuMems
	escs             escs
	flightControl    *FlightControl

	rotations     imu.Rotations
	accRotations  imu.Rotations
	gyroRotations imu.Rotations

	commandsPerSecond     int
	receiver              radioReceiver
	lastImuRead           time.Time
	imuReadInterval       time.Duration
	plotterActive         bool
	maxApplicableThrottle float64

	rollMidValue      int
	pitchlMidValue    int
	yawMidValue       int
	rotationRange     float64
	maxThrottle       float64
	minFlightThrottle float64

	plotterUdpConn      *net.UDPConn
	plotterAddress      string
	plotterDataPacket   []byte
	plotterSendBuffer   []byte
	plotterDataCounter  int
	ploterDataPerPacket int
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		startTime:             time.Now(),
		imu:                   settings.ImuMems,
		escs:                  settings.Escs,
		flightControl:         NewFlightControl(settings.Escs, settings.MinFlightThrottle, settings.PID),
		imuDataPerSecond:      settings.ImuDataPerSecond,
		receiver:              settings.Receiver,
		commandsPerSecond:     settings.CommandsPerSecond,
		lastImuRead:           time.Now(),
		imuReadInterval:       time.Second / time.Duration(settings.ImuDataPerSecond),
		plotterActive:         settings.PlotterActive,
		plotterDataPacket:     make([]byte, 0, plotter.PLOTTER_PACKET_LEN),
		plotterSendBuffer:     make([]byte, plotter.PLOTTER_PACKET_LEN),
		plotterAddress:        PLOTTER_ADDRESS,
		plotterDataCounter:    0,
		ploterDataPerPacket:   plotter.PLOTTER_DATA_PER_PACKET,
		maxApplicableThrottle: constants.MAX_APPLICABLE_THROTTLE_PERCENT,

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
	// fmt.Println("Min Flight Throttle: ", d.flightControl.minFlightThrottle)
	d.InitUdp()

	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)
	imuReadTick := utils.WithDataPerSecond(d.imuDataPerSecond)
	// escUpdates := utils.WithDataPerSecond(5)
	for running || commandOk {
		select {
		case commands, commandOk = <-commandsChannel:
			if commandOk {
				d.applyCommands(commands)
			}

		case _, running = <-ctx.Done():
			running = false
			d.plotterActive = false
			d.plotterUdpConn.Close()
		default:
			if imuReadTick.IsTime() {
				rot, err := d.imu.Read()
				if err == nil {
					d.flightControl.SetRotations(rot)
				}
			}
			// if escUpdates.IsTime() {
			// 	d.flightControl.ApplyESCThrottles()
			// }
		}
	}

	fmt.Println("Stopping Drone...")
}
