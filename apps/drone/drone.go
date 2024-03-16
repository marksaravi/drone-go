package drone

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/apps/plotter"
	"github.com/marksaravi/drone-go/devices/imu"
)

const PLOTTER_ADDRESS = "192.168.1.101:8000"

type radioReceiver interface {
	Start(ctx context.Context, wg *sync.WaitGroup, commandsPerSecond int) <-chan []byte
}

type imuMems interface {
	Read() (imu.Rotations, imu.Rotations, imu.Rotations, error)
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
	lastImuPrint          time.Time
	imuDataCounter        int
	plotterActive         bool
	maxApplicableThrottle float64

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
		flightControl:         NewFlightControl(settings.Escs),
		imuDataPerSecond:      settings.ImuDataPerSecond,
		receiver:              settings.Receiver,
		commandsPerSecond:     settings.CommandsPerSecond,
		lastImuRead:           time.Now(),
		imuReadInterval:       time.Second / time.Duration(2500),
		lastImuPrint:          time.Now(),
		imuDataCounter:        0,
		plotterActive:         settings.PlotterActive,
		rotations:             imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		accRotations:          imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyroRotations:         imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		plotterDataPacket:     make([]byte, 0, plotter.PLOTTER_PACKET_LEN),
		plotterSendBuffer:     make([]byte, plotter.PLOTTER_PACKET_LEN),
		plotterAddress:        PLOTTER_ADDRESS,
		plotterDataCounter:    0,
		ploterDataPerPacket:   plotter.PLOTTER_DATA_PER_PACKET,
		maxApplicableThrottle: 0.33,
	}
}

func (d *droneApp) readIMU() {
	if time.Since(d.lastImuRead) >= d.imuReadInterval {
		d.imuDataCounter++
		d.lastImuRead = time.Now()
		rot, acc, gyro, err := d.imu.Read()
		if err != nil {

			return
		}
		d.accRotations = acc
		d.gyroRotations = gyro
		d.rotations = rot
		if time.Since(d.lastImuPrint) >= time.Second {
			d.lastImuPrint = time.Now()
			fmt.Println(d.rotations, d.imuDataCounter)
			d.imuDataCounter = 0
		}
	}
}

func (d *droneApp) Start(ctx context.Context, wg *sync.WaitGroup) {
	var commandOk, running bool = true, true
	var commands []byte

	fmt.Println("Starting Drone...")
	d.InitUdp()

	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)

	for running || commandOk {
		d.readIMU()
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
		}
	}
	fmt.Println("Stopping Drone...")
}
