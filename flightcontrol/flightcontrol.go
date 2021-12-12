package flightcontrol

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/models"
)

type imu interface {
	ReadRotations() (models.ImuRotations, bool)
}

type pidControl interface {
	ApplyFlightCommands(flightCommands models.FlightCommands)
	ApplyRotations(rotations models.ImuRotations)
	Throttles() map[uint8]float32
}

type flightControl struct {
	pid                pidControl
	imu                imu
	throttles          chan<- models.Throttles
	onOff              chan<- bool
	escRefreshInterval time.Duration
	radio              models.Radio
	logger             chan<- models.ImuRotations
}

func NewFlightControl(
	pid pidControl,
	imu imu, throttles chan<- models.Throttles,
	onOff chan<- bool,
	escRefreshInterval time.Duration,
	radio models.Radio,
	logger chan<- models.ImuRotations,
) *flightControl {
	return &flightControl{
		pid:                pid,
		imu:                imu,
		throttles:          throttles,
		onOff:              onOff,
		escRefreshInterval: escRefreshInterval,
		radio:              radio,
		logger:             logger,
	}
}

func (fc *flightControl) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		// defer close(fc.onOff)
		// defer close(fc.throttles)
		// defer close(fc.logger)
		defer log.Println("Flight Control is stopped...")
		// fc.onOff <- true
		// time.Sleep(3 * time.Second)
		// lastEscRefresh := time.Now()
		var lastPrinted time.Time = time.Now()
		var running bool = true
		command := fc.radio.GetReceiver()
		connection := fc.radio.GetConnection()
		for running {
			// rotations, imuDataAvailable := fc.imu.ReadRotations()
			// if imuDataAvailable {
			// 	if fc.logger != nil {
			// 		fc.logger <- rotations
			// 	}
			// }
			// if time.Since(lastEscRefresh) >= fc.escRefreshInterval {
			// 	fc.pid.ApplyFlightCommands(models.FlightCommands{
			// 		Throttle: 2.5,
			// 	})
			// 	lastEscRefresh = time.Now()
			// 	fc.throttles <- fc.pid.Throttles()
			// }
			select {
			case <-ctx.Done():
				log.Println("Stopping Flight Control...")
				running = false
			case flightCommands, ok := <-command:
				if ok {
					if time.Since(lastPrinted) >= time.Second {
						showFLightCommands(flightCommands)
						lastPrinted = time.Now()
					}
				} else {
					command = nil
				}
			case connected, ok := <-connection:
				if ok {
					log.Println("connected: ", connected)
				} else {
					log.Println("channel is closed")
					connection = nil
				}
			default:
			}
		}
	}(ctx, wg)
}

func showFLightCommands(fc models.FlightCommands) {
	log.Printf("%8.2f, %8.2f, %t, %t", fc.Roll, fc.Pitch, fc.ButtonFrontLeft, fc.ButtonTopLeft)
}
