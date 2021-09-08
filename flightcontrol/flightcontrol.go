package flightcontrol

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

const (
	commandTimeCorrectionPercent float32 = 2.5
	escTimeCorrectionPercent     float32 = 3.5
	imuTimeCorrectionPercent     float32 = 8.5
)

type radio interface {
	ReceiverOn()
	ReceiveFlightData() (models.FlightCommands, bool)
	TransmitterOn()
	TransmitFlightData(models.FlightCommands) error
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
	fmt.Printf("IMU data/s: %d, ESC refresh/s: %d\n", fc.imuDataPerSecond, fc.escUpdatePerSecond)
	fc.radio.ReceiverOn()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	imuDataChannel := newImuDataChannel(ctx, &wg, fc.imu, fc.imuDataPerSecond)
	escThrottleControlChannel := newEscThrottleControlChannel(ctx, &wg, fc.esc)
	escRefresher := utils.NewTicker("esc", fc.escUpdatePerSecond, escTimeCorrectionPercent, true)
	commandChannel := newCommandChannel(ctx, &wg, fc.radio)
	fc.esc.On()
	defer fc.esc.Off()
	time.Sleep(3 * time.Second)
	go func() {
		fmt.Scanln()
		cancel()
	}()
	var throttle float32 = 0
	var running bool = true
	var rotations models.ImuRotations
	for running {
		select {
		case fd := <-commandChannel:
			throttle = fd.Throttle / 5 * 10
			if fd.IsMotorsEngaged {
				running = false
			}
		case rotations = <-imuDataChannel:
			fc.udpLogger.Send(rotations)
		case <-escRefresher:
			escThrottleControlChannel <- map[uint8]float32{
				0: throttle,
				1: throttle,
				2: throttle,
				3: throttle,
			}
		case <-ctx.Done():
			running = false
		default:
			utils.Idle()
		}
	}
	wg.Wait()
	log.Printf("stopping flightcontrol\n")
}
