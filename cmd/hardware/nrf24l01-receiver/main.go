package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	config := utils.ReadConfigs()
	_, _, radio, _ := hardware.InitDroneHardware(config)
	radio.ReceiverOn()
	var numReceive int = 0
	start := time.Now()
	for {
		if radio.IsPayloadAvailable() {
			flightdata := radio.ReceiveFlightData()
			numReceive++
			if time.Since(start) >= time.Second {
				start = time.Now()
				fmt.Println("received ", numReceive, " (", flightdata, ")")
			}
		}
	}
}
