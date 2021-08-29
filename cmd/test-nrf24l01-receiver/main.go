package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	drivers.InitHost()
	radio := utils.NewReceiverRadio()
	radio.ReceiverOn()
	var numReceive int = 0
	start := time.Now()
	for {
		if radio.IsDataAvailable() {
			flightdata := radio.ReceiveFlightData()
			numReceive++
			if time.Since(start) >= time.Second {
				start = time.Now()
				fmt.Println("received ", numReceive, " (", flightdata, ")")
			}
		}
	}
}
