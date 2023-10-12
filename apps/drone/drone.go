package drone

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/marksaravi/drone-go/apps/plotter"
	"github.com/marksaravi/drone-go/devices/imu"
)

const PLOTTER_ADDRESS = "192.168.1.101:8000"

type radioReceiver interface {
	On()
	Receive() ([]byte, bool)
}

type Imu interface {
	Read() (imu.Rotations, imu.Rotations, imu.Rotations, error)
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

func (d *droneApp) Start(ctx context.Context) {
	running := true
	d.receiver.On()
	lp := time.Now()
	d.InitUdp()
	maxCommand := time.Duration(0)
	maxIMU := time.Duration(0)
	maxUDP := time.Duration(0)
	imuPerSecond:=0
	t := byte(0)
	for running {
		select {
		default:
			ts := time.Now()
			imuok := d.ReadIMU()
			dur := time.Since(ts)
			if dur > maxIMU {
				maxIMU = dur
			}
			ts = time.Now()
			command, cmdok := d.ReceiveCommand()
			if cmdok {
				t = command[3]
			}
			dur = time.Since(ts)
			if dur > maxCommand {
				maxCommand = dur

			}
			if imuok {
				imuPerSecond++
				ts = time.Now()
				udpOk := d.SendPlotterData()
				dur = time.Since(ts)
				if udpOk && dur > maxUDP {
					maxUDP = dur
				}
			}
			if time.Since(lp) > time.Second {
				lp = time.Now()
				fmt.Printf("imu: %10v, cmd: %10v, udp: %10v, roll: %10.1f, throttle: %5d, imuPerSecond: %d\n", maxIMU, maxCommand, maxUDP, d.accRotations.Roll, t, imuPerSecond)
				maxCommand = time.Duration(0)
				maxIMU = time.Duration(0)
				maxUDP = time.Duration(0)
				imuPerSecond=0
				// if imuok {
				// 	fmt.Printf("%6.1f %6.1f %6.1f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
				// }
				// if cmd {
				// 	// fmt.Printf("%4v\n", command)
				// }
			}
		case <-ctx.Done():
			running = false
			d.plotterActive = false
			d.plotterUdpConn.Close()
		}
	}
}
