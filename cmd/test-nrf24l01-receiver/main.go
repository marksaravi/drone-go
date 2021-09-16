package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devicecreators"
	"github.com/marksaravi/drone-go/drivers"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	drivers.InitHost()
	radio := devicecreators.NewRadio()
	radio.ReceiverOn()
	var numReceive int = 0
	start := time.Now()
	var flightCommands models.FlightCommands
	for {
		fc, ia := radio.Receive()
		if ia {
			flightCommands = utils.DeserializeFlightCommand(fc)
			numReceive++
		}
		if time.Since(start) >= time.Second {
			start = time.Now()
			fmt.Println("received ", numReceive, " (", flightCommands, ")")
		}
	}
}
