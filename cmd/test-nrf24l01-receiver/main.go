package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/nrf204"
)

func main() {
	drivers.InitHost()
	radioSPIConn := drivers.NewSPIConnection(
		1,
		2,
	)
	radio := nrf204.NewNRF204("03896", "GPIO26", "-18dbm", radioSPIConn)
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
