package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/models"
)

type flightControl struct {
	imu        <-chan models.ImuRotations
	command    <-chan models.FlightCommands
	connection <-chan bool
	logger     chan<- models.ImuRotations
}

func NewFlightControl(imu <-chan models.ImuRotations, command <-chan models.FlightCommands, connection <-chan bool, logger chan<- models.ImuRotations) *flightControl {
	return &flightControl{
		imu:        imu,
		command:    command,
		connection: connection,
		logger:     logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	// var wg sync.WaitGroup
	// imuDataChannel := newImuDataChannel(ctx, &wg, fc.imu, fc.imuDataPerSecond)
	// escThrottleControlChannel := newEscThrottleControlChannel(ctx, &wg, fc.esc)
	// escRefresher := utils.NewTicker(ctx, fc.escUpdatePerSecond, 0)
	// commandChannel := newCommandChannel(ctx, &wg, fc.radio)
	// pidControl := pidcontrol.NewPIDControl()
	go func() {
		for {
			select {
			case <-fc.imu:
				// fc.logger <- rotations
			// case fc := <-fc.command:
			// 	fmt.Println(fc.ButtonFrontLeft, fc.Throttle)
			// 	// pidControl.ApplyFlightCommands(fc)
			// 	// if fc.ButtonFrontLeft {
			// 	// 	cancel()
			// 	// }
			// // case rotations := <-imuDataChannel:
			// // 	pidControl.ApplyRotations(rotations)
			// // 	fc.udpLogger.Send(rotations)
			// // case <-escRefresher:
			// // 	escThrottleControlChannel <- pidControl.Throttles()
			// case connection := <-fc.connection:
			// 	fmt.Println("connected: ", connection)
			case <-ctx.Done():
				log.Println("Flight Control Done")
				return
			}
		}
	}()
}
