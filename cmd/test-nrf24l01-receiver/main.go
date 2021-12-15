package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	hardware.InitHost()
	radioConfigs := config.ReadConfigs().FlightControl.Radio
	radio := nrf204.NewRadio(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)

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
