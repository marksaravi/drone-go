package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	config := utils.ReadConfigs()
	_, _, radio, _ := hardware.InitDroneHardware(config)
	radio.TransmitterOn()
	var roll float32 = 0
	var altitude float32 = 0
	var motorsEngaged bool = false
	var numSend int = 0
	start := time.Now()
	for range time.Tick(time.Millisecond * 20) {
		flightdata := types.FlightData{
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
