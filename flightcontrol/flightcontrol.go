package flightcontrol

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
	pidcontrol "github.com/marksaravi/drone-go/pid-control"
	"github.com/marksaravi/drone-go/utils"
)

type radio interface {
	ReceiverOn()
	Receive() ([]byte, bool)
	TransmitterOn()
	Transmit([]byte) error
}

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
	radio              radio
	udpLogger          udpLogger
}

func NewFlightControl(imuDataPerSecond int, escUpdatePerSecond int, imu imu, esc esc, radio radio, udpLogger udpLogger) *flightControl {
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
	var wg sync.WaitGroup
	fc.warmUp(cancel)
	imuDataChannel := newImuDataChannel(ctx, &wg, fc.imu, fc.imuDataPerSecond)
	escThrottleControlChannel := newEscThrottleControlChannel(ctx, &wg, fc.esc)
	escRefresher := utils.NewTicker(fc.escUpdatePerSecond, 0)
	commandChannel := newCommandChannel(ctx, &wg, fc.radio)
	pidControl := pidcontrol.NewPIDControl()
	var running bool = true
	for running {
		select {
		case fc := <-commandChannel:
			fmt.Println(fc.ButtonFrontLeft, fc.Throttle)
			pidControl.ApplyFlightCommands(fc)
			if fc.ButtonFrontLeft {
				cancel()
			}
		case rotations := <-imuDataChannel:
			pidControl.ApplyRotations(rotations)
			fc.udpLogger.Send(rotations)
		case <-escRefresher:
			escThrottleControlChannel <- pidControl.Throttles()
		case <-ctx.Done():
			running = false
		default:
			utils.Idle()
		}
	}
	wg.Wait()
	log.Printf("stopping flightcontrol\n")
}

func (fc *flightControl) warmUp(cancel context.CancelFunc) {
	fmt.Printf("IMU data/s: %d, ESC refresh/s: %d\n", fc.imuDataPerSecond, fc.escUpdatePerSecond)
	fc.radio.ReceiverOn()
	fc.esc.On()
	time.Sleep(3 * time.Second)
	go func() {
		fmt.Scanln()
		fc.esc.Off()
		cancel()
	}()
}
