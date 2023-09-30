package drone

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

const PLOTTER_BUFFER_SIZE = 9216
const PLOTTER_SAMPLE_RATE = 50

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

	plotterBuffer               []byte
	plotterDataCounter          int
	plotterSampleInterval       int
	ploterSampleIntervalCounter int
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		imu:                         settings.Imu,
		imuDataPerSecond:            settings.ImuDataPerSecond,
		receiver:                    settings.Receiver,
		commandsPerSecond:           settings.CommandsPerSecond,
		lastCommand:                 time.Now(),
		lastImuData:                 time.Now(),
		plotterActive:               settings.PlotterActive,
		rotations:                   imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		accRotations:                imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyroRotations:               imu.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		plotterBuffer:               make([]byte, 0, PLOTTER_BUFFER_SIZE),
		plotterDataCounter:          0,
		plotterSampleInterval:       settings.ImuDataPerSecond / PLOTTER_SAMPLE_RATE,
		ploterSampleIntervalCounter: 0,
	}
}

func (d *droneApp) Start(ctx context.Context) {
	running := true
	d.receiver.On()
	lp := time.Now()
	for running {
		select {
		default:
			imuok := d.ReadIMU()
			command, commandok := d.ReceiveCommand()
			if imuok {
				d.BufferPlotterData()
				d.SendPlotterData()
			}
			if (commandok) && time.Since(lp) > time.Second/10 {
				lp = time.Now()
				// if imuok {
				// 	fmt.Printf("%6.1f %6.1f %6.1f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
				// }
				if commandok {
					fmt.Printf("%4v\n", command)
				}
			}
		case <-ctx.Done():
			running = false
		}
	}
}
