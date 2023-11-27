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
	// On()
	// Receive() ([]byte, bool)
	Start(ctx context.Context, wg *sync.WaitGroup, commandsPerSecond int) <-chan []byte
}

type Imu interface {
	// Read() (imu.Rotations, imu.Rotations, imu.Rotations, error)
	Start(ctx context.Context, wg *sync.WaitGroup) <-chan imu.ImuData
}

type DroneSettings struct {
	Imu               Imu
	Receiver          radioReceiver
	ImuDataPerSecond  int
	CommandsPerSecond int
	PlotterActive     bool
}

type droneApp struct {
	startTime        time.Time
	imuDataPerSecond int
	imu              Imu

	rotations     imu.Rotations
	accRotations  imu.Rotations
	gyroRotations imu.Rotations

	commandsPerSecond int
	receiver          radioReceiver
	lastImuData       time.Time
	lastCommand       time.Time
	plotterActive     bool

	plotterUdpConn      *net.UDPConn
	plotterAddress      string
	plotterDataPacket   []byte
	plotterSendBuffer   []byte
	plotterDataCounter  int
	ploterDataPerPacket int
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		startTime:           time.Now(),
		imu:                 settings.Imu,
		imuDataPerSecond:    settings.ImuDataPerSecond,
		receiver:            settings.Receiver,
		commandsPerSecond:   settings.CommandsPerSecond,
		lastCommand:         time.Now(),
		lastImuData:         time.Now(),
		plotterActive:       settings.PlotterActive,
		rotations:           imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		accRotations:        imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyroRotations:       imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		plotterDataPacket:   make([]byte, 0, plotter.PLOTTER_PACKET_LEN),
		plotterSendBuffer:   make([]byte, plotter.PLOTTER_PACKET_LEN),
		plotterAddress:      PLOTTER_ADDRESS,
		plotterDataCounter:  0,
		ploterDataPerPacket: plotter.PLOTTER_DATA_PER_PACKET,
	}
}

func (d *droneApp) Start(ctx context.Context, wg *sync.WaitGroup) {
	var imuOk, commandOk, running bool = true, true, true

	lp := time.Now()
	lc := time.Now()
	d.InitUdp()

	commandsChannel := d.receiver.Start(ctx, wg, d.commandsPerSecond)
	imuChannel := d.imu.Start(ctx, wg)

	for running || imuOk || commandOk {
		select {
		case imuData, imuOk := <-imuChannel:
			if imuOk && imuData.Error == nil {
				d.accRotations = imuData.Accelerometer
				d.gyroRotations = imuData.Gyroscope
				d.rotations = imuData.Rotations
				if time.Since(lp) >= time.Second/2 {
					lp = time.Now()
					fmt.Println(imuData.Rotations)
				}
			}
		case command, commandOk := <-commandsChannel:
			if commandOk {
				if time.Since(lc) >= time.Second/2 {
					lc = time.Now()
					fmt.Println(command)
				}
			}

		case _, running = <-ctx.Done():
			running = false
			d.plotterActive = false
			d.plotterUdpConn.Close()
		}
	}
}
