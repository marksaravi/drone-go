package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devicecreators"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	drivers.InitHost()
	radio := devicecreators.NewRadio()
	radio.TransmitterOn()
	var roll float32 = 0
	var altitude float32 = 0
	var motorsEngaged bool = false
	var numSend int = 0
	start := time.Now()
	var id uint32 = 0
	for range time.Tick(time.Millisecond * 20) {
		flightCommands := models.FlightCommands{
			Id:       id,
			Roll:     roll,
			Pitch:    -34.53,
			Yaw:      0,
			Throttle: 13.45,
		}
		id++
		err := radio.Transmit(utils.SerializeFlightCommand(flightCommands))
		if err != nil {
			fmt.Println(err)
		}
		roll += 0.3
		altitude += 1.34
		motorsEngaged = !motorsEngaged
		numSend++
		if time.Since(start) >= time.Second {
			start = time.Now()
			fmt.Println("send ", numSend, " (", flightCommands, ")")
		}
	}
}
