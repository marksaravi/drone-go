package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

type radio interface {
	ReceiverOn()
	ReceiveFlightData() (models.FlightData, bool)
	TransmitterOn()
	TransmitFlightData(models.FlightData) error
}

func main() {
	imu, imuDataPerSecond := utils.NewImu()
	r := utils.NewRadio()
	radioChannel := NewCommandChannel(r)
	r.ReceiverOn()
	var dataPerSecond int = int(float32(imuDataPerSecond) * float32(1))
	loopDur := time.Second / time.Duration(dataPerSecond)
	const SECONDS int = 5
	var TOTAL int = SECONDS * dataPerSecond
	var counter int = 0
	fmt.Printf("Starting timer, IMU Data/Second: %d\n", imuDataPerSecond)
	start := time.Now()
	loopStart := start
	imu.ResetTime()
	var running bool = true
	for running {
		now := time.Now()
		select {
		case fd := <-radioChannel:
			fmt.Println(fd)
		default:
			r.ReceiveFlightData()
			if now.Sub(loopStart) >= loopDur {
				imu.ReadRotations()
				counter++
				if counter == TOTAL {
					running = false
				}
				loopStart = now
			}
		}
	}
	dur := time.Since(start).Seconds()
	seconds := float64(SECONDS)
	dev := (dur - seconds) * 100 / seconds
	fmt.Printf("Dur: %f, Dev: %%%f\n", dur, dev)
}

func NewCommandChannel(r radio) chan models.FlightData {
	radioChannel := make(chan models.FlightData, 10)
	go func(r radio, c chan models.FlightData) {
		ticker := time.NewTicker(time.Second / 40)
		for range ticker.C {
			if d, isOk := r.ReceiveFlightData(); isOk {
				c <- d
			}
		}
	}(r, radioChannel)
	return radioChannel
}
