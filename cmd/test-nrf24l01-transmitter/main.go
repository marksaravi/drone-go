package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/cmd/utils"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/nrf204"
)

func main() {
	drivers.InitHost()
	radio := utils.NewTransmitterRadio()
	radio.TransmitterOn()
	var roll float32 = 0
	var altitude float32 = 0
	var motorsEngaged bool = false
	var numSend int = 0
	start := time.Now()
	for range time.Tick(time.Millisecond * 20) {
		flightdata := nrf204.FlightData{
			Roll:          roll,
			Pitch:         -34.53,
			Yaw:           0,
			Throttle:      13.45,
			Altitude:      altitude,
			MotorsEngaged: motorsEngaged,
		}
		err := radio.TransmitFlightData(flightdata)
		if err != nil {
			fmt.Println(err)
		}
		roll += 0.3
		altitude += 1.34
		motorsEngaged = !motorsEngaged
		numSend++
		if time.Since(start) >= time.Second {
			start = time.Now()
			fmt.Println("send ", numSend, " (", flightdata, ")")
		}
	}
}
