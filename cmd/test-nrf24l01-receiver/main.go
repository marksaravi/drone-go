package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	drivers.InitHost()
	radio := utils.NewRadio()
	radio.ReceiverOn()
	var numReceive int = 0
	start := time.Now()
	var flightCommands models.FlightCommands
	for {
		fd, ia := radio.Receive()
		if ia {
			flightCommands = fd
			numReceive++
		}
		if time.Since(start) >= time.Second {
			start = time.Now()
			fmt.Println("received ", numReceive, " (", flightCommands, ")")
		}
	}
}
