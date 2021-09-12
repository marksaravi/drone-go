package flightcontrol

import (
	"context"
	"log"
	"sync"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

func newEscThrottleControlChannel(ctx context.Context, wg *sync.WaitGroup, escdevice esc) chan<- map[uint8]float32 {
	wg.Add(1)
	escChannel := make(chan map[uint8]float32, 10)
	go func(escdev esc, ch chan map[uint8]float32) {
		defer wg.Done()
		var throttles map[uint8]float32
		for {
			select {
			case throttles = <-ch:
				escdev.SetThrottles(throttles)
			case <-ctx.Done():
				log.Printf("stoping esc channel\n")
				return
			default:
				utils.Idle()
			}
		}
	}(escdevice, escChannel)
	return escChannel
}

func newImuDataChannel(ctx context.Context, wg *sync.WaitGroup, imudev imu, dataPerSecond int) <-chan models.ImuRotations {
	wg.Add(1)
	imuDataChannel := make(chan models.ImuRotations, 10)
	go func(imudev imu, ch chan models.ImuRotations) {
		defer wg.Done()
		ticker := utils.NewTicker(dataPerSecond, 0)
		for range ticker {
			ch <- imudev.ReadRotations()
			select {
			case <-ctx.Done():
				log.Printf("stoping imu channel\n")
				return
			default:
			}
		}
	}(imudev, imuDataChannel)
	return imuDataChannel
}

func newCommandChannel(ctx context.Context, wg *sync.WaitGroup, r radio) <-chan models.FlightCommands {
	wg.Add(1)
	radioChannel := make(chan models.FlightCommands, 10)
	go func(r radio, c chan models.FlightCommands) {
		defer wg.Done()
		ticker := utils.NewTicker(40, 0)
		for {
			select {
			case <-ctx.Done():
				log.Printf("stoping command channel\n")
				return
			case <-ticker:
				if d, isOk := r.Receive(); isOk {
					c <- utils.DeserializeFlightCommand(d)
				}
			default:
				utils.Idle()
			}
		}
	}(r, radioChannel)
	return radioChannel
}