package drone

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

type radioReceiver interface {
	On()
	Receive() ([]byte, bool)
}

type Imu interface {
	Read() (imu.Rotations, error)
}

type DroneSettings struct {
	Imu               Imu
	Receiver          radioReceiver
	ImuDataPerSecond  int
	CommandsPerSecond int
	PlotterActive     bool
}

type droneApp struct {
	imuDataPerSecond  int
	commandsPerSecond int
	imu               Imu
	receiver          radioReceiver
	lastImuData       time.Time
	lastCommand       time.Time
	plotterActive     bool
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		imu:               settings.Imu,
		imuDataPerSecond:  settings.ImuDataPerSecond,
		receiver:          settings.Receiver,
		commandsPerSecond: settings.CommandsPerSecond,
		lastCommand:       time.Now(),
		lastImuData:       time.Now(),
		plotterActive:     settings.PlotterActive,
	}
}

func (d *droneApp) Start(ctx context.Context) {
	running := true
	d.receiver.On()
	lp := time.Now()
	for running {
		select {
		default:
			rotations, imuok := d.ReadIMU()
			command, commandok := d.ReceiveCommand()
			if imuok {
				d.PlotterData(rotations)
			}
			if (imuok || commandok) && time.Since(lp) > time.Second/10 {
				lp = time.Now()
				if imuok {
					fmt.Printf("%6.1f %6.1f %6.1f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
				}
				if commandok {
					fmt.Printf("%4v\n", command)
				}
			}
		case <-ctx.Done():
			running = false
		}
	}
}
