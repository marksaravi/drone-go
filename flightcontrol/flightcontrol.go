package flightcontrol

import (
	"context"
	"fmt"
	"sync"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type flightControl struct {
	command    chan models.FlightCommands
	connection chan bool
}

func NewFlightControl(command chan models.FlightCommands, connection chan bool) *flightControl {
	return &flightControl{
		command:    command,
		connection: connection,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// var wg sync.WaitGroup
	// imuDataChannel := newImuDataChannel(ctx, &wg, fc.imu, fc.imuDataPerSecond)
	// escThrottleControlChannel := newEscThrottleControlChannel(ctx, &wg, fc.esc)
	// escRefresher := utils.NewTicker(ctx, fc.escUpdatePerSecond, 0)
	// commandChannel := newCommandChannel(ctx, &wg, fc.radio)
	// pidControl := pidcontrol.NewPIDControl()
	var running bool = true
	for running {
		select {
		case fc := <-fc.command:
			fmt.Println(fc.ButtonFrontLeft, fc.Throttle)
			// pidControl.ApplyFlightCommands(fc)
			// if fc.ButtonFrontLeft {
			// 	cancel()
			// }
		// case rotations := <-imuDataChannel:
		// 	pidControl.ApplyRotations(rotations)
		// 	fc.udpLogger.Send(rotations)
		// case <-escRefresher:
		// 	escThrottleControlChannel <- pidControl.Throttles()
		case connection := <-fc.connection:
			fmt.Println("connected: ", connection)
		case <-ctx.Done():
			running = false
		default:
			utils.Idle()
		}
	}
}
