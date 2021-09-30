package flightcontrol

import (
	"context"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/devices/radioreceiver"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type imu interface {
	ReadRotations() models.ImuRotations
	ResetTime()
}
type esc interface {
	Off()
	On()
	SetThrottles(map[uint8]float32)
}

type udpLogger interface {
	Send(models.ImuRotations)
}

type flightControl struct {
	imuDataPerSecond   int
	escUpdatePerSecond int
	imu                imu
	esc                esc
	radio              models.RadioLink
	udpLogger          udpLogger
}

func NewFlightControl(imuDataPerSecond int, escUpdatePerSecond int, imu imu, esc esc, radio models.RadioLink, udpLogger udpLogger) *flightControl {
	return &flightControl{
		imuDataPerSecond:   imuDataPerSecond,
		escUpdatePerSecond: escUpdatePerSecond,
		imu:                imu,
		esc:                esc,
		radio:              radio,
		udpLogger:          udpLogger,
	}
}

func (fc *flightControl) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("Press ENTER to exit.")
		fmt.Scanln()
		cancel()
	}()
	// var wg sync.WaitGroup
	// imuDataChannel := newImuDataChannel(ctx, &wg, fc.imu, fc.imuDataPerSecond)
	// escThrottleControlChannel := newEscThrottleControlChannel(ctx, &wg, fc.esc)
	// escRefresher := utils.NewTicker(ctx, fc.escUpdatePerSecond, 0)
	// commandChannel := newCommandChannel(ctx, &wg, fc.radio)
	// pidControl := pidcontrol.NewPIDControl()
	const heatbeatPerSecond int = 4
	const commandPerSecond int = 20
	receiver := radioreceiver.NewRadioReceiver(ctx, fc.radio, commandPerSecond, heatbeatPerSecond)
	var running bool = true
	for running {
		select {
		case fc := <-receiver.Command:
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
		case connection := <-receiver.Connection:
			fmt.Println("connected: ", connection)
		case <-ctx.Done():
			running = false
		default:
			utils.Idle()
		}
	}
	// wg.Wait()
	log.Printf("stopping flightcontrol\n")
}
