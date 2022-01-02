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
	var roll float32 = 0
	var altitude float32 = 0
	var motorsEngaged bool = false
	var numSend int = 0
	start := time.Now()
	var id uint32 = 0
	for range time.Tick(time.Millisecond * 20) {
		flightCommands := models.FlightCommands{
			Roll:     14,
			Pitch:    15,
			Yaw:      0,
			Throttle: 17,
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
