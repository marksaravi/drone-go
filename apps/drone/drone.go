package drone

import (
	"context"
	"fmt"
	"time"
)

type radioReceiver interface {
	On()
	Receive() ([]byte, bool)
}

type DroneSettings struct {
	Receiver          radioReceiver
	CommandsPerSecond int
}

type droneApp struct {
	receiver          radioReceiver
	commandsPerSecond int
	lastCommand       time.Time
}

func NewDrone(settings DroneSettings) *droneApp {
	return &droneApp{
		receiver:          settings.Receiver,
		commandsPerSecond: settings.CommandsPerSecond,
		lastCommand:       time.Now(),
	}
}

func (d *droneApp) Start(ctx context.Context) {
	running := false
	for running {
		select {
		default:
			command, ok := d.ReceiveCommand()
			if ok {
				fmt.Println(command)
			}
		case <-ctx.Done():
			running = false
		}
	}
}
